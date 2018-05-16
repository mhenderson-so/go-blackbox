package blackbox

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/openpgp"
)

// AdminAdd adds an administrator to a blackbox repository. It searches for
// gpg public keys by email address.
func AdminAdd(email, directory string) ([]string, error) {
	email = strings.TrimSpace(email)

	//Check if the user is currently an admin. If they are, skip them
	admins, err := ListAdmins()
	if err != nil {
		return nil, err
	}
	found := false
	for _, admin := range admins {
		if strings.ToLower(admin) == strings.ToLower(email) {
			found = true
		}
	}

	if !found {
		admins = append(admins, email)
		sort.Strings(admins)
	}

	//Open the users personal gpg keyring
	personalKeyringPath := ""

	//Which keyring should we look in? If path is blank, use the users home
	//keyring. If it's not blank, use the provided path to look for the keyring.
	if directory != "" {
		personalKeyringPath = filepath.Join(directory, "pubring.gpg")
	} else {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		homeDir := usr.HomeDir
		personalKeyringPath = filepath.Join(homeDir, ".gnupg", "pubring.gpg")
	}

	//Read the users public keyring and get the list of entities in the keyring
	keyringFileBuffer, err := os.Open(personalKeyringPath)
	if err != nil {
		return nil, err
	}
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}

	//Read the blackbox public keyring so we can add the new keys to it if we
	//find any to add
	blackboxFileBuffer, err := os.OpenFile(PublicKeyringPath(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer blackboxFileBuffer.Close()
	blackboxEntityList, err := openpgp.ReadKeyRing(blackboxFileBuffer)
	if err != nil {
		return nil, err
	}
	//Go through the users public keys and find ones that match the provided email
	//address. If we find one, check if it's currently in the blackbox keyring.
	//If it isn't, add the description of the found key to the output.
	var addedKeys []string
	var foundMatch bool
	for _, entity := range entityList { //Our users public keyrings
		var keyAdded bool
		for _, identity := range entity.Identities { //The identities contained within this keyring
			if identity.UserId.Email == email { //We have a matching identity
				foundMatch = true
				//Check if we have a key that matches this ID already
				matchingKeys := blackboxEntityList.KeysById(entity.PrimaryKey.KeyId)
				//If we have zero matching keys, then this is a key to be added
				if len(matchingKeys) == 0 {
					//Add the entity to the list
					blackboxEntityList = append(blackboxEntityList, entity)

					//We've added the key, this is just to generate a nice string to pass
					//up the stack so the user can see which keys have been added. This could be
					//a struct with these fields which might be nicer in the long run.
					fingerprint := hex.EncodeToString(entity.PrimaryKey.Fingerprint[:])
					fingerprint = fingerprint[len(fingerprint)-8 : len(fingerprint)]
					thisDesc := fmt.Sprintf("%s: %s <%s>", fingerprint, identity.UserId.Name, identity.UserId.Email)
					if identity.UserId.Comment != "" {
						thisDesc = fmt.Sprintf("%s (%s)", thisDesc, identity.UserId.Comment)
					}
					addedKeys = append(addedKeys, thisDesc)
					keyAdded = true
				}
			}
		}

		if keyAdded {
			//We found a key to add, so put it onto the end of the keyring. The file seek() position will
			//already be at the end of the file due to ReadKeyRing() earlier, so we don't need to move around.
			err := entity.Serialize(blackboxFileBuffer)
			if err != nil {
				return nil, err
			}

			//And write out the new admins file
			adminFile := filepath.Join(blackboxData, blackboxAdminsFile)
			adminData := []byte(strings.Join(admins, "\n"))
			ioutil.WriteFile(adminFile, adminData, 644)
			continue
		}
	}

	//If we went through all of this but we never found a matching key, return an error saying so
	if !foundMatch {
		return nil, fmt.Errorf("Unable to find key matching %s", email)
	}

	return addedKeys, nil
}

// AdminCleanup compares the list of blackbox admin users against what is
// present in the keychain, and removes users from the keychain who should
// not be present
func AdminCleanup() ([]string, error) {
	bbAdmins, err := ListAdmins()
	if err != nil {
		return nil, err
	}

	blackboxFileBuffer, err := os.OpenFile(PublicKeyringPath(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer blackboxFileBuffer.Close()
	blackboxEntityList, err := openpgp.ReadKeyRing(blackboxFileBuffer)
	if err != nil {
		return nil, err
	}

	var removedAdmins []string
	var finalEntityList openpgp.EntityList
	for _, entity := range blackboxEntityList {
		var foundMatch bool
		for _, identity := range entity.Identities {
			for _, admin := range bbAdmins {
				if !foundMatch && strings.TrimSpace(strings.ToLower(admin)) == strings.TrimSpace(strings.ToLower(identity.UserId.Email)) {
					//Check if expired, remove if true
					//Commented out for now, as doing this may require a run to update blackbox-admins.txt if we remove
					//all the keys for a given user.
					//if identity.SelfSignature.KeyExpired(time.Now()) {
					//	continue
					//}
					foundMatch = true
					finalEntityList = append(finalEntityList, entity)
					continue
				}
			}
			if !foundMatch {
				fingerprint := hex.EncodeToString(entity.PrimaryKey.Fingerprint[:])
				fingerprint = fingerprint[len(fingerprint)-8 : len(fingerprint)]
				thisDesc := fmt.Sprintf("%s: %s <%s>", fingerprint, identity.UserId.Name, identity.UserId.Email)
				if identity.UserId.Comment != "" {
					thisDesc = fmt.Sprintf("%s (%s)", thisDesc, identity.UserId.Comment)
				}

				removedAdmins = append(removedAdmins, thisDesc)
			}
		}
	}

	blackboxFileBuffer.Truncate(0)
	blackboxFileBuffer.Seek(0, io.SeekStart)
	for _, ident := range finalEntityList {
		ident.Serialize(blackboxFileBuffer)
	}

	return removedAdmins, nil
}

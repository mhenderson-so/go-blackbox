package blackbox

import (
	"encoding/hex"
	"fmt"
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
	blackboxFileBuffer, err := os.Open(PublicKeyringPath())
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
	for _, entity := range entityList { //Our users public keyrings
		var keyAdded bool
		for _, identity := range entity.Identities { //The identities contained within this keyring
			if identity.UserId.Email == email { //We have a matching identity

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
				}
			}
		}
		if keyAdded {
			continue
		}
	}

	return addedKeys, nil
}

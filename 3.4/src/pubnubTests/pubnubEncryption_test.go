package pubnubTests

import (
	"testing"
	"pubnubMessaging"
	"fmt"
)

func TestYayDecryptionBasic(t *testing.T) {
 	message := "q/xJqqN6qbiZMXYmiQC1Fw==";
    //decrypt
    decrypted := pubnubMessaging.DecryptString("enigma", message);

    if("yay!" == decrypted){    	
    	fmt.Println("Yay decryption passed.") 
    } else {
    	t.Error("Yay decryption failed.");
    }
}

func TestYayEncryptionBasic(t *testing.T) {
 	message := "yay!";
    //decrypt
    encrypted := pubnubMessaging.EncryptString("enigma", message);

    if("q/xJqqN6qbiZMXYmiQC1Fw==" == encrypted){
    	fmt.Println("Yay encryption passed.") 
    } else {
    	t.Error("Yay encryption failed.");
    }
}
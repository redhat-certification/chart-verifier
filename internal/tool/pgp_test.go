package tool

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	keyfileName    = "../../tests/charts/psql-service/0.1.11/psql-service-0.1.11.tgz.key"
	expectedDigest = "1cc31121e86388fad29e4cc6fc6660f102f43d8c52ce5f7d54e134c3cb94adc2"
)

func TestPGPKeyEncoding(t *testing.T) {
	encodedKey, encodeErr := GetEncodedKey(keyfileName)
	require.NoError(t, encodeErr)
	require.True(t, len(encodedKey) > 0)

	keyDigest, digestErr := GetPublicKeyDigest(encodedKey)
	require.NoError(t, digestErr)
	require.Equal(t, expectedDigest, keyDigest)

	decodedKey, decodeErr := GetDecodedKey(encodedKey)
	require.NoError(t, decodeErr)
	require.True(t, len(decodedKey) > 0)

	keyBytes, readErr := ioutil.ReadFile(keyfileName)
	require.NoError(t, readErr)
	require.Equal(t, keyBytes, decodedKey)
}

func TestDigest2(t *testing.T) {
	base64Cmd := exec.Command("base64", "-i", keyfileName)
	base64KeyFromCmd, _ := base64Cmd.Output()

	base64Key := strings.Trim(string(base64KeyFromCmd), " -\n")
	base64Key = strings.Replace(base64Key, "\n", "", -1)

	encodedKey, encodeErr := GetEncodedKey(keyfileName)
	require.NoError(t, encodeErr)
	require.True(t, len(encodedKey) > 0)

	require.Equal(t, encodedKey, base64Key)

	base64Echo := exec.Command("echo", base64Key)
	sha256cmd := exec.Command("sha256sum")

	sha256Value := bytes.NewBufferString("")
	sha256cmd.Stdin, _ = base64Echo.StdoutPipe()
	sha256cmd.Stdout = sha256Value

	_ = sha256cmd.Start()
	_ = base64Echo.Run()
	_ = sha256cmd.Wait()

	shaResponseSplit := strings.Split(sha256Value.String(), " ")
	require.Equal(t, expectedDigest, strings.TrimRight(shaResponseSplit[0], " -\n"))
}

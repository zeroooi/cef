package aegis

import "testing"

func TestClient_GetConfig(t *testing.T) {
	cli := NewAegisClient("http://0.0.0.0:11470/")
	resp, err := cli.GetConfig("browser-config", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

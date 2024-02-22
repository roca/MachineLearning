package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type MyTokenCredential struct{}

func (m *MyTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		//ExpiresOn: time.Now().Add(1 * time.Hour),
	}, nil
}

func azure() {
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	username := os.Getenv("SVC_ACCOUNT")
	//password := os.Getenv("SVC_PASSWORD")
	//version := os.Getenv("OPENAI_API_VERSION")
	scope := os.Getenv("AZURE_APPLICATION_SCOPE")

	//tenant := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)

	scopes := []string{scope}

	token, err := getToken(clientID, tenantID, scopes)
	if err != nil {
		log.Fatalf("ERROR getToken: %s", err)
	}
	fmt.Println("TOKEN:", token)

	key, err := getkey(token)
	if err != nil {
		log.Fatalf("ERROR getkey: %s", err)
	}
	fmt.Println("OPENAPI_KEY:", key)

	client, err := getClient(key)
	if err != nil {
		log.Fatalf("ERROR getClient: %s", err)
	}
	fmt.Println("Client established.")

	_ = client

	// azlog.SetListener(func(event azlog.Event, s string) {
	// 	fmt.Println("azidentity: ", s)
	// })

	modelDeploymentID := "RegnADA002"
	resp, err := client.GetEmbeddings(context.Background(), azopenai.EmbeddingsOptions{
		Input:          []string{"The food was delicious and the waiter..."},
		DeploymentName: &modelDeploymentID,
		User: &username,
	}, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	fmt.Println("Got embeddings for input: ", resp)

	for _, embed := range resp.Data {
		// embed.Embedding contains the embeddings for this input index.
		fmt.Fprintf(os.Stderr, "Got embeddings for input %d\n", *embed.Index)
	}

}

func getToken(clientID, tenantID string, scopes []string) (string, error) {
	o := azidentity.InteractiveBrowserCredentialOptions{
		//AdditionallyAllowedTenants: []string{"*"},
		ClientID: clientID,
		TenantID: tenantID,
	}

	//dac, err := azidentity.NewUsernamePasswordCredential(tenantID, clientID, username, password, &o)
	dac, err := azidentity.NewInteractiveBrowserCredential(&o)
	if err != nil {
		return "", err
	}

	token, err := dac.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes:   scopes,
		TenantID: tenantID,
	})
	return token.Token, nil
}

func getkey(jwtToken string) (string, error) {
	endpoint := os.Getenv("AZURE_ENDPOINT")

	body := []byte(`{"key":"openaikey2"}`)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("client: could not create request: %s\n", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

	client := http.Client{
		Timeout: 1 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client: error making http request: %s\n", err)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
		return "", fmt.Errorf("Error while reading the response bytes: %s", err)
	}

	fmt.Println(string([]byte(body)))
	return string([]byte(body)), nil
}

func getClient(key string) (*azopenai.Client, error) {
	//endpoint := os.Getenv("AZURE_ENDPOINT")
	//baseEndpoint := os.Getenv("OPENAI_API_BASE")
	azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	openaiKey := azcore.NewKeyCredential(key)

	options := azopenai.ClientOptions{}

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, openaiKey, &options) //
	if err != nil {
		log.Fatalf("ERROR NewClient: %s", err)
	}
	return client, nil
}

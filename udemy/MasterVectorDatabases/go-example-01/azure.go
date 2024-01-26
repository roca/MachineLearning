package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type MyTokenCredential struct{}

func (m *MyTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IjVCM25SeHRRN2ppOGVORGMzRnkwNUtmOTdaRSJ9.eyJhdWQiOiI5NTUzNzNiYy1kODBmLTRiMzUtYTc5Mi0yMmE4NTM3YzA1OWIiLCJpc3MiOiJodHRwczovL2xvZ2luLm1pY3Jvc29mdG9ubGluZS5jb20vM2U5YWFkZjgtNmExNi00OTBmLThkY2QtYzY4ODYwY2FhZTBiL3YyLjAiLCJpYXQiOjE3MDYyODQyNDYsIm5iZiI6MTcwNjI4NDI0NiwiZXhwIjoxNzA2Mjg5MDI1LCJhaW8iOiJBVVFBdS84VkFBQUFTNEYvajNIY2pBWVJzYVM2eXlldVE4NHVidWFOcWY1N0Q1MUdXMU9SbHBncEszVUNtME04NDNuMUtmeENHSjl2ZXNHZ3grVXlVMlVsTGNmR1FhTjl5dz09IiwiYXpwIjoiOTU1MzczYmMtZDgwZi00YjM1LWE3OTItMjJhODUzN2MwNTliIiwiYXpwYWNyIjoiMCIsIm5hbWUiOiJzdmMgY2hhdGdwdC1hcGkiLCJvaWQiOiIzZjM1MzJkMi1jZjE1LTQ1NDktYWRiNy0yYWRhYmNlYmUwNTgiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJzdmMuY2hhdGdwdC1hcGlAcmVnZW5lcm9uLmNvbSIsInJoIjoiMC5BWFlBLUsyYVBoWnFEMG1OemNhSVlNcXVDN3h6VTVVUDJEVkxwNUlpcUZOOEJadTBBTDQuIiwic2NwIjoiQXBpVXNlci5SZWFkIiwic3ViIjoiVlhOQkQ4Y3VMZmhJV0dsanRQX1JyQk5IUWpaLXlvdHZ1cENhQkR5TzBoVSIsInRpZCI6IjNlOWFhZGY4LTZhMTYtNDkwZi04ZGNkLWM2ODg2MGNhYWUwYiIsInV0aSI6IjFGUE8tYUt2dGsteFczOHg4N0RsQUEiLCJ2ZXIiOiIyLjAifQ.BJKvcRN_6dkesY-hEwJcHLMy-6xIk3PFWX0tP8buog5edgJITB7PYwsGINcNSgXWIA1TeyUJZB3RAior1Olqh-sp1Pv9m0x0vCpKzkqv7GexlCOQK_MfJJ5dpBSgZ0n7nxrnKyrmIz7rH4tKkm0j-HJ5Kz6XvwgzLUAdj83sEM51yqAMT2JjU-UrjE7POWzK8c5Lz_58cmugmtZ-lmW-kwbR1Z07pkhREYnLrX3LpHjOVcIGaNKgRGTfOZuU_jB-CyRIxUiN-QO-RIjTF3xWwUa2ZD0B8QhUewgfGZ6KfZojbT5fqAxBcNRV5euuJmWhPDL8MmLse0-H51BdXM7byw",
		//ExpiresOn: time.Now().Add(1 * time.Hour),
	}, nil
}

func azure() {
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	endpoint := os.Getenv("AZURE_ENDPOINT")
	//baseEndpoint := os.Getenv("OPENAI_API_BASE")
	username := os.Getenv("SVC_ACCOUNT")
	password := os.Getenv("SVC_PASSWORD")
	//version := os.Getenv("OPENAI_API_VERSION")
	scope := os.Getenv("AZURE_APPLICATION_SCOPE")
	//azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

	//tenant := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)
	scopes := []string{scope}

	dac, err := azidentity.NewUsernamePasswordCredential(tenantID, clientID, username, password, nil)
	if err != nil {
		log.Fatal(err)
	}
	_ = dac

	token, _ := dac.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes:   scopes,
		TenantID: tenantID,
	})
	fmt.Println("TOKEN:", token.Token)

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClient(endpoint, dac, nil) // NewClientForOpenAI(endpoint, key, nil) //
	if err != nil {
		log.Fatalf("ERROR NewClient: %s", err)
	}
	_ = client
	modelDeploymentID := "RegnADA002"
	resp, err := client.GetEmbeddings(context.TODO(), azopenai.EmbeddingsOptions{
		Input:          []string{"The food was delicious and the waiter..."},
		DeploymentName: &modelDeploymentID,
	}, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, embed := range resp.Data {
		// embed.Embedding contains the embeddings for this input index.
		fmt.Fprintf(os.Stderr, "Got embeddings for input %d\n", *embed.Index)
	}

}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
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
	endpoint := os.Getenv("AZURE_ENDPOINT")
	//baseEndpoint := os.Getenv("OPENAI_API_BASE")
	username := os.Getenv("SVC_ACCOUNT")
	password := os.Getenv("SVC_PASSWORD")
	//version := os.Getenv("OPENAI_API_VERSION")
	scope := os.Getenv("AZURE_APPLICATION_SCOPE")
	//azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

	tenant := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)
	scopes := []string{scope}

	httpClient := &http.Client{}

	clientApp, err := public.New(clientID, public.WithAuthority(tenant), public.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalln("Fatal error calling public.New: ", err)
	}

	acquire_tokens_result2, err := clientApp.AcquireTokenByUsernamePassword(context.Background(), scopes, username, password)
	if err != nil {
		log.Fatalln("Fatal error calling clientApp.AcquireTokenByUsernamePassword", err)
	}
	fmt.Println("TOKEN:", acquire_tokens_result2.AccessToken)

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

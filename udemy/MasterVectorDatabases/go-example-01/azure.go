package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
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
	//endpoint := os.Getenv("AZURE_ENDPOINT")
	//baseEndpoint := os.Getenv("OPENAI_API_BASE")
	username := os.Getenv("SVC_ACCOUNT")
	//password := os.Getenv("SVC_PASSWORD")
	//version := os.Getenv("OPENAI_API_VERSION")
	scope := os.Getenv("AZURE_APPLICATION_SCOPE")
	azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

	//tenant := fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)

	scopes := []string{scope}

	//azlog.SetEvents(azidentity.EventAuthentication)

	// o := azidentity.UsernamePasswordCredentialOptions{
	// 	ClientOptions: azcore.ClientOptions{
	// 		//APIVersion: "2023-03-15-preview",
	// 		Logging: policy.LogOptions{
	// 			IncludeBody: true,
	// 		},
	// 		Retry: policy.RetryOptions{
	// 			MaxRetries: 10,
	// 		},

	// 	},
	// }

	o := azidentity.InteractiveBrowserCredentialOptions{
		//AdditionallyAllowedTenants: []string{"*"},
		ClientID: clientID,
		TenantID: tenantID,
	}

	//dac, err := azidentity.NewUsernamePasswordCredential(tenantID, clientID, username, password, &o)
	dac, err := azidentity.NewInteractiveBrowserCredential(&o)
	if err != nil {
		log.Fatalf("Error calling NewUsernamePasswordCredential: %v", err)
	}

	token, err := dac.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes:   scopes,
		TenantID: tenantID,
	})
	if err != nil {
		log.Fatalf("Error calling GetToken: %v", err)
	}
	fmt.Println("TOKEN:", token)

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	co := azopenai.ClientOptions{}
	client, err := azopenai.NewClient(azureOpenAIEndpoint, dac, &co) // NewClientForOpenAI(endpoint, key, nil) //
	if err != nil {
		log.Fatalf("ERROR NewClient: %s", err)
	}
	fmt.Println("Client established.")

	azlog.SetListener(func(event azlog.Event, s string) {
		fmt.Println("azidentity: ", s)
	})

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

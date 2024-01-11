package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func azure() {
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	username := os.Getenv("SVC_ACCOUNT")
	password := os.Getenv("SVC_PASSWORD")
	version := os.Getenv("OPENAI_API_VERSION")
	//scope := os.Getenv("AZURE_APPLICATION_SCOPE")

	//azureOpenAIKey := os.Getenv("API_API_KEY")
	modelDeploymentID := "text-embedding-ada-002"

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := endpoint

	// if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
	// 	fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
	// 	return
	// }

	dac, err := azidentity.NewUsernamePasswordCredential(tenantID, clientID, username, password, &azidentity.UsernamePasswordCredentialOptions{
		ClientOptions: policy.ClientOptions{
			APIVersion: version,
			
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// token, _ := dac.GetToken(context.Background(), policy.TokenRequestOptions{
	// 	TenantID: tenantID,
	// 	Scopes:   []string{scope},
	// })

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClient(azureOpenAIEndpoint, dac, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

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

#  Code ref: https://ritscm.regeneron.com/projects/DSEAIM/repos/secureaccess/browse/securekeyapp.py

# pip install --upgrade pip
# pip install openai
# pip install msal
# pip install retry

import os
import openai
import requests
import simplejson as json
from retry import retry
from msal import PublicClientApplication, ConfidentialClientApplication, ClientApplication
import time

@retry (tries=3, delay=2)
def getapikey():
    # These are the AZURE parameters needed by the client application
    # In this scheme, the user does not have to worry about the api key, all of that is handled
    # at the AZURE Back End.   The API key is rotated hourly.
    # We can handle these as Kuberbetes secrets, or as environment variables.
    client_id = os.getenv("AZURE_CLIENT_ID")
    tenant_id = os.getenv("AZURE_TENANT_ID")
    endpoint = os.getenv("AZURE_ENDPOINT")

    scopes = [os.getenv("AZURE_APPLICATION_SCOPE")]

    app = ClientApplication(
        client_id=client_id,
        authority="https://login.microsoftonline.com/" + tenant_id
    )


    acquire_tokens_result = app.acquire_token_by_username_password(username=os.getenv("SVC_ACCOUNT"),
                                                                   password=os.getenv("SVC_PASSWORD"),
                                                                   scopes=scopes)
    if 'error' in acquire_tokens_result:
        print("Error: " + acquire_tokens_result['error'])
        print("Description: " + acquire_tokens_result['error_description'])
        return 2
    else:
        header_token = {"Authorization": "Bearer {}".format(acquire_tokens_result['access_token'])}
        rt = requests.post(url=endpoint, headers=header_token, data=b'{"key":"openaikey2"}')
        return rt.json()

openai.api_type = os.getenv("OPENAI_API_TYPE")
openai.api_base = os.getenv("OPENAI_API_BASE")
openai.api_version = os.getenv("OPENAI_API_VERSION")
# We are dynamically getting the key from AZURE>   Access is based on the service account/ad group combination
openai.api_key = getapikey()
print(openai.api_key)

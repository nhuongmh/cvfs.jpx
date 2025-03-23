import os.path

from google.auth.transport.requests import Request
from google.oauth2.credentials import Credentials
from google.oauth2 import service_account
from google_auth_oauthlib.flow import InstalledAppFlow
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError

# If modifying these scopes, delete the file token.json.

class GoogleSheet:
    def __init__(self, creds_data, spreadsheet_id, sheet_name):
        self.creds_data = creds_data
        self.spreadsheet_id = spreadsheet_id
        self.sheet_name = sheet_name

    def load(self):
      """Shows basic usage of the Sheets API.
      Prints values from a sample spreadsheet.
      """
      # The file token.json stores the user's access and refresh tokens, and is
      # created automatically when the authorization flow completes for the first
      # time.
      # Load the base64 encoded credentials from the .secret.env file
      creds = service_account.Credentials.from_service_account_info(
            self.creds_data, scopes=["https://www.googleapis.com/auth/spreadsheets"]
      )

      try:
        service = build("sheets", "v4", credentials=creds)

        # Call the Sheets API
        sheet = service.spreadsheets()
        result = (
            sheet.values()
            .get(spreadsheetId=self.spreadsheet_id, range=f'{self.sheet_name}!A2:E')
            .execute()
        )
        values = result.get("values", [])

        if not values:
          print("No data found.")
          return

        return values
      except HttpError as err:
        print(err)

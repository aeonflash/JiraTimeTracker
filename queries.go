package main

const CurrentUserQuery = `query currentUser {
  me {
    user {
      ... on AtlassianAccountUser {
        accountId
        accountStatus
        name
        picture
      }
    }
  }
}`
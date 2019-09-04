# Blog module

## Requirements

This module defines the required components for blog.

- A blog is where a user posts their article
- Every user can post article and has permission delete only their article
- Users can like and comment to articles

### State

- #### User

  - ID
  - Username
  - Bio

- #### Blog

  - ID
  - Owner
  - Title
  - Description
  - Articles
  - CreatedAt

- #### Article

  - ID
  - Title
  - Content
  - CreatedAt
  - DeleteAt

### Messages

- #### Create User

  - Username
  - Bio

- #### Create Blog

  - Title
  - Description

- #### Create Article

  - Title
  - Content
  - DeleteAt

# Louie's Library

A site for sharing kindle books with friends.

Currently in alpha, these are the implimented features:

  - Books
    - Uploaded in .mobi format, data fetched from book.google api
  - Requests
  - Collections
    - Record your reading progress
  - Invite System for new users
  - User messages/inbox (early stages)
  - Extra Utilities:
    - Youtube Video/Playlist to MP3 (useful for audiobooks)

TODO:
  - Collections
    - Set year of completion
    - Reading goal
  - Site Announcements
  - Mobi to EPub/PDF Conversion
  - How-To Video for Download/Install on Kindle
  - Consider document-based DB (mongo?)


Implimented using Go 1.14 and PostgreSQL 12. Hosted on Linode with Ubuntu 20.04LTS.
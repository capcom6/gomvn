# Changelog

All notable user-visible changes to GoMVN are documented in this file.

## [0.4.10] - 2026-07-22

### New Features

- **AWS credential-chain authentication** — S3 storage can now use the AWS SDK default credential chain when `login` and `password` are omitted from the GoMVN configuration

### Bug Fixes

- **Complete S3 listings** — Repository and artifact browsing now follows every S3 list-objects page instead of stopping after the first 1,000 matching keys

Notes

URL: `[scheme:][//[userinfo@]host][/]path[?query][#fragment]`

Example:
 <https://username:password@github.com:12345/bomctl/bomctl.git#sbom.cdx.json?ref=main>

- scheme: https
- userinfo: [username, password]
- host: github.com:12345 (parse port with url.Port())
- path: bomctl/bomctl.git
- fragment: sbom.cdx.json
- Query: ref=main (parse query with url.Query() and get value with query.Get("ref"))

# Existing

| client | fetch  | push   | init   | required |
|:-------|:-------|:-------|:-------|:---------|
| git    | "fragment" | "scheme", "username", "password", "hostname", "fragment" | "scheme", "username", "password", "hostname", "port", "path", "gitRef", "fragment" | "scheme", "hostname", "path", "gitRef", "fragment" |
| github | "username", "password", "path",  | N/A | "scheme", "username", "password", "hostname", "port", "path" | "scheme", "hostname", "path" |
| gitlab | "path", "fragment" | N/A | "scheme", "hostname", "port", "path", "fragment" | "scheme", "hostname", "path", "fragment" |
| http   | "username", "password", "hostname" | N/A | "scheme", "username", "password", "hostname", "port", "path", "gitRef", "fragment" | "scheme", "hostname" |
| oci    | "hostname", "port", "path", "tag" | "username", "password", "hostname", "port", "path", "tag" | "scheme", "username", "password", "hostname", "port", "path", "tag", "digest" | "scheme", "hostname", "path" |

# Proposed

| client | fetch  | push   | init   | required |
|:-------|:-------|:-------|:-------|:---------|
| git    | "fragment" | "scheme", "username", "password", "hostname", "fragment" | "scheme", "username", "password", "hostname", "port", "path", "gitRef (query)", "fragment" | "scheme", "hostname", "path", "gitRef (query)", "fragment" |
| github | "username", "password", "path",  | N/A | "scheme", "username", "password", "hostname", "port", "path" | "scheme", "hostname", "path" |
| gitlab | "path", "gitref (query)" | N/A | "scheme", "hostname", "port", "path", "gitref (query)" | "scheme", "hostname", "path", "gitref (query)" |
| http   | "username", "password", "hostname" | N/A | "scheme", "username", "password", "hostname", "port", "path", "query", "fragment" | "scheme", "hostname" |
| oci    | "hostname", "port", "path", "gitref (query)" | "username", "password", "hostname", "port", "path", "gitref (query)" | "scheme", "username", "password", "hostname", "port", "path", "gitref (query)" | "scheme", "hostname", "path" |

| scheme  | host | path     | fragment | Gitref | Client | Notes |
|:-------:|:----:|:--------:|:--------:|:------:|:------:|:-----:|
| git     | +    | +        | *        | *      | git    | client should validate all required fields are present? |
| github  | +    | +        | *        | *      | github | client should validate all required fields are present? |
| gitlab  | +    | +        | *        | *      | gitlab | client should validate all required fields are present? |
| http    | +    | +        | *        | *      | http   | client should validate all required fields are present? |
| oci     | +    | +        | *        | *      | oci    | client should validate all required fields are present? |
| http(s) | +    | +        | +        | *      | git    | Contains fragment, only git client needs that |
| http(s) | +    | `.git`   | *        | *      | git    | path contains `.git`, send to git client (optional) |
| http(s) | +    | `gitlab` | -        | -      | git    | path contains `gitlab` no fragment and no gitref, send to gitlab client |
| http(s) | +    | `github` | -        | +      | git    | path contains `github` no fragment and a gitref, send to github client |

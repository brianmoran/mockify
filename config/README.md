# Configuration

## Explanation

The configuration file is either YAML or JSON, and contains an array of route objects that Mockify will respond to.

Each route must have 3 fields:

   - `route` [`string`]: the URI with any variables portions of it defined by `{}` 
   - `methods` [`array`]: REST method(s) that this route will provide responses for. Must contain at least one
   - `responses` [`array`]: an array of responses that each have 6 fields:
      - `uri` [`string`]: the URI that this reponses is used for (no variable parts here)
      - `method` [`string`]: REST method for the response
      - _(Optional)_ `requestBody` [`string`]: If any part of the request body matches this string, that response will be used (only used if it has a value, can't be an empty string. Have the highest matching priority)
      - _(Optional)_ `requestHeader` [`string`]: Must be in the format `Key: Value` (i.e. `Authorization: foo-bar`). If `Key` is found in the request, and `Key`'s value is `Value` that response will be used (only used if it has a value, can't be an empty string. Have the second highest matching priority)
      - `statusCode` [`integer`]: Response status code
      - `headers` [`object`]: Response headers
      - `body` [`object`]: Response body
         - If `headers` contains `"Content-Type": "application/json"` the response body will be converted to JSON, otherwise the response is sent as a string but the the content-type provided so and `application/xml` can be interpreted correctly.

### Matching priorities

There are three different priorities:

1. The routes with the `requestBody` has the highest matching priority. This means that they will be chosen first as the response when Mockify is called
1. The routes with the `requestHeader` has the second highest matching priority. This means that they will be chosen if a route with the highest priority hasn't been found
1. The routes without `requestBody` nor `requestHeader` has the lowest matching priority 

## Examples

See the configuration files in this folder.

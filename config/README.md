# Configuration (examples and explanations)

The configuration file is YAML and is an array of route objects that mochky will respond to.  

Each routes has 3 fields: 
   - route: the basic URI with any variables portions of it defined by `{}`. 
   - methods: REST methods that this route will provide responses for.
   - responses is an array of responses that each have 6 fields:
      - uri: the URI that this reponses is used for (no variable parts here)
      - method: REST method oin the request
      - requestBody: any part of the requeast body (only used if it had a value other that empty string.
      - statusCode: Response status code
      - headers: Response headers
      - body: Response body (if the headers contain the return type of `"Content-Type": "application/json"` the response body will be converted to JSON, otherwise the response is sent as a string but the the content-type provided so and `application/xml` can be intepreted correcttly.
      
Notes on Response matching:
    - if the requestBody is absent or is an empty string the matching will be made only on the URI and Method otherwise all 3 fields are used to match the proper response.
    
 
    
    
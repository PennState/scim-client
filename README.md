# scim-client - a SCIM [1] client for the Go [2] programming language

## Features

- Schema [3]
  - Standard resources
    - ``User`` struct
    - ``EnterpriseUser`` struct
  - Custom resources
    - Provides the ``Resource`` struct for creation of custom resources
    - Provides the ``Extension`` interface for the creation of custom resource extensions
- Protocol [4]
  - ``ErrorResponse`` struct
  - ``ListResponse`` struct

## Planned features

- Schema
  - Standard resources
    - ``Group`` struct
  - Server resources
    - ``ResourceType`` struct
    - ``Schema`` struct
    - ``ServiceProviderConfig`` struct
- Protocol
- Client
  - Operations
    - HTTP "CRUD" operations (POST, GET, PUT, DELETE)
    - Partial update (PATCH)
    - Filtered search (GET)
  - Conflict management
    - eTag support
    - Last modified support
  - Authentication modes
    - Unauthenticated
    - Basic authentication
    - OAuth2 Client Credential authentication
    - Custom authentication
- Validation
  - Validate resources against a server-provided schema
  - Validate resources against a custom schema
  - Validate resources against the default schema (``User``, ``EnterpriseUser`` and ``Group`` only)

## User notes

### Client creation and initialization

### Resource manipulation

Once a SCIM client has been created, retrieving a SCIM resource can be easily
accomplished as shown in the following code example:

```go
func changeManager(client scim.Client, employeeID string, managerID string) error {
    e, e_err := client.Get("User", employeeID) // [1]
    if e_err != nil {
        return e_err
    }

    m, m_err := client.Get("User", managerID) // [1]
    if m_err != nil {
        return m_err
    }

    var manager Manager
    manager.Value = m.ID
    manager.Reference = "../Users/" + m.ID
    manager.DisplayName = m.DisplayName

    var enterpriseUser EnterpriseUser
    err := e.GetExtension(&enterpriseUser) // [2]
    if eu_err != nil {
        return eu_err
    }
    eu.Manager = manager
    err = e.UpdateExtension(&enterpriseUser) // [3]

    return client.Update("User", e) // [4]
}
```

This code illustrates how to:

1. Retrieve a resource from a server
2. Get one of the resource's extensions
3. Update the local representation of the resource
4. Update a resource's representation on the server

### Custom ``Resource`` creation

A custom resource can be created by creating a ``struct`` annotated for
JSON marshaling/unmarshaling and adding the ``scim.Resource`` field as
shown in the following example:

```go
//Organization represents some hierarchy of an arbitrary organization
//including (URI) an optional reference to a parent organization as well
//as to zero or more child organizations.  As with other SCIM references
//the Parent and Children references may be absolute or relative URIs.
type Organization struct {
    scim.Resource
    Name     string                `json:"name"`      //Name is the organization's name - e.g. "Tour Promotion"
    Type     string                `json:"type"`      //Type is the organization's type - e.g. "Department"
    Parent   OrganizationReference `json:"$parent"`   //Parent is a URI reference to a parent organization
    Children OrganizationReference `json:"$children"` //Children is a URI reference to zero or more child organizations
}

//OrganizationReference is a string containing an absolute or relative
//URL to another Organization.
type OrganizationReference string
```

The scim client library will marshal to and unmarshal from the following
JSON:

```json
{
    "id": "430beb5c-a361-4c04-b308-2845789a496e",
    "schemas": ["urn:com:example:2.0:Organization"],
    "name": "Tour Promotion",
    "type": "Department",
    "parent": "../Organizations/4a7741a3-a436-4a52-a6d5-149e6c1b9578",
    "children": [
        "../Organizations/7eb59c46-35a4-4443-b8c1-5de8be88f973",
        "../Organizations/66506f29-8c44-414e-b52d-a993b94f370c",
        "../Organizations/0a365d4f-10e5-45c5-ae05-ee5184b59627"
    ],
    "meta": {
        "resourceType": "Organization",
        "created": "2010-01-23T04:56:22Z",
        "lastModified": "2011-05-13T04:42:34Z",
        "version": "W/3694e05e9dff590",
        "location":"https://example.com/v2/Organizations/430beb5c-a361-4c04-b308-2845789a496e"
    }
}
```

### Custom ``Extension`` creation

Custom SCIM ``Extension``s can also be created using a techique similar
to that shown above.  In this case, a SCIM ``Extension`` is identified
by implementing the ``GetUrn() string`` function as shown in the
following example:

```go
//Pantry is a SCIM Extension that allows employees and employers to track
//the amount owed to the pantry (or to the employee).
type Pantry struct {
    Building string  `json:"building"` //Building is the location housing the employee's office
    Office   string  `json:"office"`   //Office is the room number of the employee's office
    Balance  float64 `json:"balance"`  //Balance is the amount the employee owes to the pantry (if negative).  Credits can be represented by positive Balance values
}

//GetURN returns the SCIM Extension's URN (identifier) and, more importantly
//identifies the Pantry struct as a SCIM extension.
func (p Pantry) GetURN() string {
    return "urn:com:example:2.0:Pantry"
}
```

When extending a SCIM ``User`` resource, the ``Extension`` will marshall to
and unmarshal from the following JSON:

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:User",
        "urn:com:example:2.0:Pantry"
    ],
    "id": "2819c223-7f76-453a-919d-413861904646",
    "externalId": "701984",
    "userName": "bjensen@example.com",
    "urn:com:example:2.0:Pantry": {
        "building": "Technology Support Building",
        "office": "202BB",
        "balance": 4.25
    },
    "meta": {
        "resourceType": "User",
        "created": "2010-01-23T04:56:22Z",
        "lastModified": "2011-05-13T04:42:34Z",
        "version": "W/a330bc54f0671c9",
        "location": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
    }
}
```

## Developer notes

This project is a Go module and therefore does not need to be located
in the ``$GOPATH`` hierarchy.  Use the ``./...`` construct with your Go
commands as follows:

```text
go build ./...                        // builds the project's packages
go test ./...                         // runs the tests in all the project's packages
go test -cover ./...                  // runs the tests in all the project's packages showing coverage information
go test -coverprofile cover.out ./... // runs the tests in all the project's packages showing coverage information and generating a coverage report
go tool cover -html=cover.out         // displays the generated coverage report in your browser of choice
```

## References

1. [System for Cross-domain Identity Management (SCIM)](http://www.simplecloud.info)
2. [Go programming language](https://golang.org)
3. [RFC7643 - SCIM: Core Schema](https://tools.ietf.org/html/rfc7643)
4. [RFC7644 - SCIM: Protocol](https://tools.ietf.org/html/rfc7644)

# HAL

![Build](https://github.com/read-with-cortex/hal/actions/workflows/build-and-test.yml/badge.svg)
[![GoDoc](https://godoc.org/github.com/read-with-cortex/hal?status.svg)](https://godoc.org/github.com/read-with-cortex/hal)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://pmoule.mit-license.org)

A **Go** implementation of **Hypertext Application Language (HAL)**.
It provides essential data structures and features a JSON generator
to produce **JSON** output as proposed in [JSON Hypertext Application Language](https://tools.ietf.org/html/draft-kelly-json-hal).
As an add-on it also supports the [The HAL-FORMS Media Type](https://rwcbook.github.io/hal-forms/) based on the <span style="color:red">**Working Draft**</span> from 2021-03-03.

## Features
- HAL API
    - Create root **Resource Object**.
    - Supports "_links" property.

        Define **Link Relations** and assign **Link Object** value(s).

    - Supports "_embedded" property.

        Define **Link Relations** and assign **Resource Object** value(s).
    - Supports "curies".

        Define CURIE **Link Objects** and assign to defined **Link Relations**.
- JSON generator to produce HAL Document
- Tools to simplify HAL document creation
- HAL-FORMS
  - JSON generator to produce HAL-FORMS documents
  - Supports creating relations in HAL pointing to HAL-FORMS documents

## Usage
### Preliminary stuff
Download and install `go2hal` into your GOPATH.
```
go get github.com/read-with-cortex/hal
```
Import the `hal` package to get started.
```go
import "github.com/read-with-cortex/hal/hal"
```
### Create a Resource and JSON generation

First create a Resource Object as your HAL document's root element.

```go
root := hal.NewResourceObject()
```

This is all you need, to create a valid HAL document.

Next, create an `Encoder` and call it's `Encode` function to generate valid JSON.

```go
encoder := hal.NewJSONEncoder()
bytes, error := encoder.Encode(root)
```
Generated JSON
```JSON
{}
```

There's potential for more :smile:

### Add a Link Relation
So let's add a `self` Link Relation. Additionally we attach a single Link Object.
```go
link := &hal.LinkObject{ Href: "/docwhoapi/doctors"}

self, _ := hal.NewLinkRelation("self") //skipped error handling
self.SetLink(link)

root.AddLink(self)
```
Since `self` is a IANA registered link relation name, **go2hal** provides a shortcut
```go
self := hal.NewSelfLinkRelation()
```

Generated JSON
```JSON
{
    "_links": {
        "self": {
            "href": "/docwhoapi/doctors"
        }
    }
}
```
### Resource state
To add some resource state, you can add some properties.
These properties must be valid JSON

(Currently, this is not checked)
```go
data := root.Data()
data["doctorCount"] = 12
```
Generated JSON
```JSON
{
    "_links": {
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "doctorCount": 12
}
```
It might be a bit cumbersome to manually map all properties from
any of your DTOs. There is a more convenient way.
```go
// a simple struct
type DoctorsInfo struct {
    Content     string  `json:"content"`
    DoctorCount int     `json:"doctorCount"` 
    From        string  `json:"from"`
    Until       string  `json:"until"`
}

info := DoctorsInfo{12, "All actors of the Doctor.", "1963", "today"}
root.AddData(info)
```
Generated JSON
```JSON
{
    "_links": {
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "content": "All actors of the Doctor."
    "doctorCount": 12,
    "from": "1963"
    "until": "today"
}
```
Both ways of adding state can be combined. But already existing properties are replaced.
### Embedding Resources
Now, let's embed some resources.
```go
// a simple struct for actors of the doctor
type Actor struct {
    ID   int    `json:"-"`
    Name string `json:"name"`
}

actors := []Actor {
    Actor{1, "William Hartnell"},
    Actor{2, "Patrick Troughton"},
}

// convert the actors to resources
var embeddedActors []hal.Resource

for _, actor := range actors {
    href := fmt.Sprintf("/docwhoapi/doctors/%d", actor.ID)
    selfLink, _ := hal.NewLinkObject(href)

    self, _ := hal.NewLinkRelation("self")
    self.SetLink(selfLink)

    embeddedActor := hal.NewResourceObject()
    embeddedActor.AddLink(self)
    embeddedActor.AddData(actor)
    embeddedActors = append(embeddedActors, embeddedActor)
}

doctors, _ := hal.NewResourceRelation("doctors")
doctors.SetResources(embeddedActors)

root.AddResource(doctors)
```
Generated JSON
```JSON
{
    "_links": {
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "_embedded": {
        "doctors": [
            {
                "_links": {
                    "self": {
                        "href": "/docwhoapi/doctors/1"
                    }
                },
                "name": "William Hartnell"
            },
            {
                "_links": {
                    "self": {
                        "href": "/docwhoapi/doctors/2"
                    }
                },
                "name": "Patrick Troughton"
            }
        ]
    },
    "content": "All actors of the Doctor."
    "doctorCount": 12,
    "from": "1963"
    "until": "today"
}
```
### CURIEs
A Resource Object can have a set of CURIE links. Same for used Link Relations, that are capable of setting a CURIE link.
```go
curieLink, _ := hal.NewCurieLink("doc", "http://example.com/docs/relations/{rel}")
curieLinks := []*hal.LinkObject {curieLink}
root.AddCurieLinks(curieLinks)

doctors.SetCurieLink(curieLink)
```
Generated JSON
```JSON
{
    "_links": {
        "curies": [
            {
                "href": "http://example.com/docs/relations/{rel}",
                "templated": true,
                "name": "doc"
            }
        ],
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "_embedded": {
            "doc:doctors": [
                {
                    "_links": {
                        "self": {
                            "href": "/docwhoapi/doctors/1"
                        }
                    },
                    "name": "William Hartnell"
                },
                {
                    "_links": {
                        "self": {
                            "href": "/docwhoapi/doctors/2"
                        }
                    },
                    "name": "Patrick Troughton"
                }
            ]
        },
    "doctorCount": 12
}
```
### Relations and the array vs single value discussion
I'm aware of existing discussions regarding Relations and the type of assigned values.
I simply deal with this topic by leaving the decision to the developer whether to
assign a single value or an array value.

**go2hal** provides a `LinkRelation` and a `ResourceRelation`. Both are capable of holding a single value or an array value by providing special setter functions.

**Example**
```go
link := hal.NewLinkObject("myHref")
linkRelation := hal.NewLinkRelation("linkRel")

resource := hal.NewResourceObject()
resource.Data()["value"] = "myValue"
resourceRelation := hal.NewResourceRelation("resourceRel")
```
Assign single values
```go
linkRelation.SetLink(link)
resourceRelation.SetResource(resource)
```
Generated JSON
```JSON
{
    "_links": {
        "linkRel": {
            "href": "myHref"
        }
    },
    "_embedded": {
            "resourceRel": {
                "value": "myValue"
            }
        },
}
```
Assign array values
```go
linkRelation.SetLinks([]*LinkObject{link})
resourceRelation.SetResources([]hal.Resource{resource})
```
Generated JSON
```JSON
{
    "_links": {
        "linkRel": [
            {
                "href": "myHref"
            }
        ]
    },
    "_embedded": {
        "resourceRel": [
            {
                "value": "myValue"
            }
        ]
    },
}
```
**go2hal** does not evaluate the assigned values. The developer is fully responsible of this.

CURIEs are an exception to this. As stated in the [specification](https://tools.ietf.org/html/draft-kelly-json-hal#section-8.2)
CURIEs are always an array of Link Objects.

### Tooling
To simplify creating HAL documents, **go2hal** provides a `ResourceFactory`.
Initialise it with a set of CURIE links.
```go
curieLink, _ := hal.NewCurieLink("doc", "http://example.com/docs/relations/{rel}")
curieLinks := []*hal.LinkObject {curieLink}

factory := NewResourceFactory(curieLinks)
```
Create a root resource with it's `self` link.
```go
self := "/docwhoapi/doctors"
root := factory.CreateRootResource(self)
```
Generated JSON
```JSON
{
    "_links": {
        "curies": [
            {
                "href": "http://example.com/docs/relations/{rel}",
                "templated": true,
                "name": "doc"
            }
        ],
        "self": {
            "href": "/docwhoapi/doctors"
        }
    }
}
```
Create a link and an embedded resource and it's `ResourceRelation`.
```go
//assign a CURIE link by name
companions := factory.CreateLink("companions", "/docwhoapi/companions", "doc")
root.AddLink(companions)

doctor := Actor{1, "William Hartnell"}
embeddedSelf := fmt.Sprintf("/docwhoapi/doctors/%d", doctor.ID)
embeddedDoctor := factory.CreateEmbeddedResource(embeddedSelf)
embeddedDoctor.AddData(doctor)

//assign a CURIE link by name
doctorLink := factory.CreateResourceLink("hartnell", "doc")
doctorLink.SetResource(embeddedDoctor)

root.AddResource(doctorLink)
```
Generated JSON
```JSON
{
    "_links": {
        "curies": [
            {
                "href": "http://example.com/docs/relations/{rel}",
                "templated": true,
                "name": "doc"
            }
        ],
        "doc:companions": {
            "href": "/docwhoapi/companions"
        },
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "_embedded": {
        "doc:hartnell": [
            {
                "_links": {
                    "self": {
                        "href": "/docwhoapi/doctors/1"
                    }
                },
                "name": "William Hartnell"
            }
        ]
    },
}
```
### HAL-FORMS
Let's create a link relation pointing to a **HAL-FORMS** document for creating a new resource.
```go
createLink, _ := halforms.NewHALFormsRelation{"create", "/docwhoapi/hal-forms/create-doctor"} //skipped error handling
root.AddLink(createLink)
```
Generated JSON
```JSON
{
    "_links": {
        "create": {
            "href": "/docwhoapi/hal-forms/create-doctor",
            "type": "application/prs.hal-forms+json"
        }
    }
}
```

A **HAL-FORMS** documents can be created and encoded to **JSON** this way:

```go
href := "www.example.com"
document := NewDocument(href)
encoder := halforms.NewJSONEncoder()
bytes, _ := encoder.Encode(document) // skipped error handling
```

Generated JSON
```JSON
{
    "_links": {
        "self": {
            "href": "www.example.com",
            "type": "application/prs.hal-forms+json"
        }
    },
    "_templates": {}
}
```

The created document can be enriched with detailed form filed information.
```go
// create a new template describing how to POST a new value to target URI
template := NewTemplate()
template.Method = http.MethodPost
template.Target = "/docwhoapi/doctors"
// create a name property for the form.
property := NewProperty("name")
property.Prompt = "Name"
property.Placeholder = "the doctor's name"
property.Required = true

template.Properties = append(template.Properties, property)
document.AddTemplate(template)
```
Generated JSON
```JSON
{
    "_links": {
        "self": {
            "href": "/docwhoapi/hal-forms/create-doctor",
            "type": "application/prs.hal-forms+json"
        }
    },
    "_templates": {
        "default": {
            "contentType": "application/json",
            "key": "default",
            "method": "POST",
            "properties": [
                {
                    "name": "name",
                    "prompt": "Name",
                    "readOnly": false,
                    "regex": "",
                    "required": true,
                    "templated": false,
                    "value": "",
                    "cols": 40,
                    "max": "",
                    "maxLength": "",
                    "min": "",
                    "minLength": "",
                    "placeholder": "the doctor's name",
                    "rows": 5,
                    "step": 0,
                    "type": "text"
                }
            ],
            "target": "/docwhoapi/doctors",
            "title": "default"
        }
    }
}
```
## Documentation

See package documentation:

[![GoDoc](https://godoc.org/github.com/read-with-cortex/hal?status.svg)](https://godoc.org/github.com/read-with-cortex/hal)

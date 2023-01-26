<!--
 Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Lesser General Public License v3.0 as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Lesser General Public License v3.0 for more details.
 You should have received a copy of the GNU Lesser General Public License v3.0
 along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
-->
# Ipaas Framework

Integration platform as a service (iPaaS) is a set of automated tools that integrate software applications that are deployed in different environments. Large businesses that run enterprise-level systems often use iPaaS to integrate applications and data that live on premises and in both public and private clouds

It is a main folder called ipass and contains subfolders as app,files,model,repository and utils

## App
App folder contains handler and route files

## Files

Files folder contains subfolders as apis,features and mappers

### Api's
API stands for “Application Programming Interface.” An API is a software intermediary that allows two applications to talk to each other. In other words, an API is the messenger that delivers your request to the provider that you're requesting it from and then delivers the response back to you

#### Sample 1
 ```json
{
    "name": "task_name",
    "description": "task description",
    "endpoint": "{host}/api/v1/sample",
    "scheme": "https",
    "session_variables": [
        {
            "key": "total_batch",
            "value": "mention the featureSessionVariable Key which contain the array response to execute batch_request",
            "type": "variable"
        },
        {
            "key": "per_batch",
            "value": "5",
            "type": "static"
        },
        {
            "key": "batch_size",
            "value": "",
            "type": "dynamic",
            "props": [
                {
                    "name": "GetLength",
                    "params": [
                        {
                            "key": "array",
                            "value": "mention the featureSessionVariable Key which contain the array response to execute batch_request",
                            "type": "variable"
                        }
                    ]
                }
            ]
        },
        {
            "key": "current_batch",
            "value": "",
            "type": "static"
        }
    ],
    "request_configuration": {
        "method": "GET",
        "type": "batch",
        "request_details": [
            {
                "key": "mention what all are fields need to be update each time",
                "value": "",
                "type": "dynamic",
                "props":[]
            }
        ]
    },
    "completion_norms": [],
    "success_callback": [],
    "payload": {
        "type": "body",
        "headers": [],
        "body": [],
        "query_params": [],
        "oauth": []
    },
    "response": {
        "session_variables": []
    }
}
```
##### Available Sample 1 fields are:

- **name** (*str*, required)
    - The human-readable name of the task
- **description** (*str*)
    - Extended description for the task, in reStructuredText
- **endpoint** (*str*)
    - An API Endpoint is the URL for a server or a service
- **scheme** (*str*)
    - Schemes define which transfer protocols you want your API to use
- **session variables** (*list(obj)*)
    - Session variables are special variables that exist only while the user's session with your application is active. Session variables are specific to each visitor to your site. They are used to store user-specific information that needs to be accessed by multiple pages in a web application
    
 | session_variables | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| key | *str* | batch_size | 
| value | *str* | Mark_123123dsf | 
| type | *str* | dynamic |
| props | *list(obj)* | {"name": "GetLength","params": [{"key": "array","value": "mention the featureSessionVariable Key which contain the array response to execute batch_request","type": "variable"}]} |

- **request configuration** ( *obj*)
   - Request Config. These are the available config options for making requests

 | request_configuration | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| method | *str* | GET |
| type | *str* | batch |
| request_details | *list(obj)* | { "key":"mention what all are fields need to be update each time" "value": "" "type": "dynamic","props":[]} |

- **completion_norms** (*list(obj)*)
- **success_callback** (*list(obj)*)
- **payload** (*obj*)
    - The body of your request and response message

 | payload | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| type | *str* | body |
| headers | *list(dict)* | [ "Content-Type":"application/json"] |
| body | *list(obj)* | [{"name":"task"}] |
| query_params | *list(dict)* | [{"query":"task_name"}] |
| oauth | *list(dict)* | [{"token":"token"}] |


- **response** (*obj*)
    - A response is a message the server receives in return for a Request we send

 | response | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| session_variables | *list(obj)* | [{"session_key:"Mark_123123dsf"}] |

#### Sample 2

```json
{
    "name": "task_name",
    "description": "task description",
    "endpoint": "{host}/api/v1/sample",
    "scheme": "http",
    "request_configuration": {
        "method": "GET",
        "type": "paginated",
        "request_details": []
    },
    "completion_norms": [
        {
            "name": "function_name",
            "params": []
        }
    ],
    "success_callback": [
        {
            "name": "function_name",
            "params": []
        }
    ],
    "session_variables": [
        {
            "key": "",
            "value": "",
            "type": "static or variable or dynamic",
            "props": [
                {
                    "name": "function_name",
                    "params": []
                }
            ]
        }
    ],
    "payload": {
        "type": "body",
        "headers": [],
        "body": [],
        "query_params": [],
        "oauth": []
    },
    "response": {
        "session_variables": []
    }
}
```
##### Available Sample 2 fields are:

- **name** (*str*, required)
    - The human-readable name of the task
- **description** (*str*)
    - Extended description for the task, in reStructuredText
- **endpoint** (*str*)
    - An API Endpoint is the URL for a server or a service
- **scheme** (*str*)
    - Schemes define which transfer protocols you want your API to use
- **request configuration** ( *obj*)
    - Request Config. These are the available config options for making requests

 | request_configuration | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| method | *str* | GET |
| type | *str* | paginated |
| request_details | *list(obj)* | { "key":"mention what all are fields need to be update each time" "value": "" "type": "dynamic","props":[]} |

- **completion_norms** (*list(obj)*)

 | completion_norms | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| name | *str* | function_name |
| params | *list(obj)* | [{"key": "value"}] |

- **success_callback** (*list(obj)*)

 | success_callback | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| name | *str* | function_name |
| params | *list(obj)* | [{"key": "value"}] |

- **session variables** (*list(obj)*)
    - Session variables are special variables that exist only while the user's session with your application is active. Session variables are specific to each visitor to your site. They are used to store user-specific information that needs to be accessed by multiple pages in a web application
    
 | session_variables | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| key | *str* | per_batch |
| value | *str* | 5 |
| type | *str* | static or variable or dynamic | 
| props | *list(obj)* | { "name": "function_name",
"params": []} |

- **payload** (*obj*)
    - The body of your request and response message

 | payload | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| type | *str* | body |
| headers | *list(dict)* | [ "Content-Type":"application/json"] |
| body | *list(obj)* | [{"name":"task"}] |
| query_params | *list(dict)* | [{"query":"task_name"}] |
| oauth | *list(dict)* | [{"token":"token"}] |

- **response** (*obj*)
    - A response is a message the server receives in return for a Request we send

 | response | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| session_variables | *list(obj)* | [{"session_key:"Mark_123123dsf"}] |
#### Sample 3

```json
{
    "name": "task_name",
    "description": "task description",
    "endpoint": "{host}/api/v1/sample",
    "scheme": "http",
    "request_configuration": {
        "method": "GET",
        "type": "single",
        "request_details": []
    },
    "completion_norms": [],
    "success_callback":[
        {
            "name": "function_name",
            "params": []
        }
    ],
    "session_variables": [
        {
            "key": "",
            "value": "",
            "type": "static or variable or dynamic",
            "props": [
                {
                    "name": "function_name",
                    "params": []
                }
            ]
        }
    ],
    "payload": {
        "type": "body",
        "headers": [],
        "body": [],
        "query_params": [],
        "oauth":[]
    },
    "response": {
        "session_variables": []
    }
}

```
##### Available Sample 3 fields are:

- **name** (*str*, required)
    - The human-readable name of the task
- **description** (*str*)
    - Extended description for the task, in reStructuredText
- **endpoint** (*str*)
    - An API Endpoint is the URL for a server or a service
- **scheme** (*str*)
    - Schemes define which transfer protocols you want your API to use
- **request configuration** ( *obj*)
    - Request Config. These are the available config options for making requests

 | request_configuration | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| method | *str* | GET |
| type | *str* | single |
| request_details | *list(obj)* | { "key":"mention what all are fields need to be update each time" "value": "" "type": "dynamic","props":[]} |

- **completion_norms** (*list(obj)*)

- **success_callback** (*list(obj)*)

 | success_callback | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| name | *str* | function_name |
| params | *list(obj)* | [{"key": "value"}] |

- **session variables** (*list(obj)*)
    - Session variables are special variables that exist only while the user's session with your application is active. Session variables are specific to each visitor to your site. They are used to store user-specific information that needs to be accessed by multiple pages in a web application
    
 | session_variables | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| key | *str* | per_batch |
| value | *str* | 5 |
| type | *str* | static or variable or dynamic | 
| props | *list(obj)* | { "name": "function_name",
"params": []} |

- **payload** (*obj*)
    - The body of your request and response message

 | payload | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| type | *str* | body |
| headers | *list(dict)* | [ "Content-Type":"application/json"] |
| body | *list(obj)* | [{"name":"task"}] |
| query_params | *list(dict)* | [{"query":"task_name"}] |
| oauth | *list(dict)* | [{"token":"token"}] |

- **response** (*obj*)
    - A response is a message the server receives in return for a Request we send

 | response | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| session_variables | *list(obj)* | [{"session_key:"Mark_123123dsf"}] |
### Features

The feature folder serves to write the organisation api's and Synchronise the data from the required modules

#### Sample 1

```json
{
    "name": "template_name",
    "description": "template description",
    "on_failure":"patch (or) rollback",
    "tasks": [
      {
        "name": "task_name_1",
        "link": "/template/files/apis/task_name_1.json",
        "type": "endpoint",
        "plugin_auth_required": false,
        "dependents":[],
        "on_success":[],
        "on_success_skip_tasks":[],
        "priority": 1
      },
      {
        "name": "task_name_2",
        "link": "/template/files/apis/task_name_2.json",
        "type": "endpoint",
        "plugin_auth_required": false,
        "dependents":["task_name_1"],
        "on_success":["task_name_3"],
        "on_success_skip_tasks":["task_name_3"],
        "mapper": {
          "name": "woocommerce.get_products_mapper",
          "link": "./integrations/woocommerce/files/apis/woocommerce.get_products_mapper.json",
          "input_type": "array",
          "input_key":"data.products"
        },
        "priority": 2
      },
      {
        "name": "task_name_3",
        "link": "/template/files/apis/task_name_3.json",
        "type": "endpoint",
        "plugin_auth_required": false,
        "dependents":[],
        "on_success":[],
        "on_success_skip_tasks":[],
        "required_mapped_output": [
          { "name": "woocommerce.get_products", "type": "payload"}
        ],
        "priority": 3
      }
    ]
  }
```
##### Available Sample 1 fields are:

- **name** (*str*, required)
    - The human-readable name of the app
- **description** (*str*)
    - Extended description for the app, in reStructuredText
- **on_failure** (*str*)
    - on failure task it will be patch or rollback
- **tasks** (*list(obj)*)
    - It's a necessary step on the road towards project completion
- **example for task 1**

 | tasks | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| name | *str* | task_name_1 |
| link | *str* | /template/files/apis/task_name_1.json |
| type | *str* | endpoint | 
| plugin_auth_required | bool | false |
| dependents | *list(obj)* | [{"dependents":"values"}] | 
| on_success | *list(obj)* | [{"on_success":"values"}] | 
| on_success_skip_tasks | *list(obj)* | [{"on_success_skip_tasks":"values"}] |
| priority | *int* | 1 | 

- **example for task 2**

 | tasks | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| name | *str* | task_name_2 |
| link | *str* | /template/files/apis/task_name_2.json |
| type | *str* | endpoint | 
| plugin_auth_required | bool | false |
| dependents | *list(obj)* | task_name_1 | 
| on_success | *list(obj)* | task_name_3 | 
| on_success_skip_tasks | *list(obj)* | task_name_3 |
| mapper | *obj* | { "name": "woocommerce get_products_mapper","link": "./internal/app integrations/woocommerce/files/apis/woocommerce get_products_mapper.json","input_type": "array" "input_key":"data.products"} |
| priority | *int* | 2 | 

- **example for task 3**

 | tasks | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| name | *str* | task_name_3 |
| link | *str* | /template/files/apis/task_name_3.json |
| type | *str* | endpoint | 
| plugin_auth_required | bool | false |
| dependents | *list(obj)* | [{"dependents":"values"}] | 
| on_success | *list(obj)* | [{"on_success":"values"}] | 
| on_success_skip_tasks | *list(obj)* | [{"on_success_skip_tasks":"values"}] |
| required_mapped_output | *list(obj)* | { "name":"woocommerce.get_products", "type":"payload"} |
| priority | *int* | 3 | 

### Mappers

Mapper is a Go package that allows for mapping multiple prepared SQL statements with name "keys". This can be used to manage resources and close all statements when your program closes
#### Sample 1

```json
{
  "name": "woocommerce.get_products_mapper",
  "description": "maps products from woocommerce to eunimart",
  "fields": [
    {
      "id": "field_sku_code",
      "scope": "local",
      "input": "sku",
      "input_type": "keyvalue",
      "output": "sku_code",
      "output_type": "keyvalue",
      "default_value": "",
      "execution_type": "",
      "helper_function": []
    },
    {
      "id": "field_product_name",
      "scope": "local",
      "input": "name",
      "input_type": "keyvalue",
      "output": "product_name",
      "output_type": "keyvalue",
      "default_value": "",
      "execution_type": "",
      "helper_function": []
    },
    {
      "id": "field_package_height",
      "scope": "local",
      "input": "height",
      "input_type": "keyvalue",
      "output": "package_height",
      "output_type": "keyvalue",
      "default_value": "",
      "fields": [],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "field_package_length",
      "scope": "local",
      "input": "length",
      "input_type": "keyvalue",
      "output": "package_length",
      "output_type": "keyvalue",
      "default_value": "",
      "fields": [],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "field_package_width",
      "scope": "local",
      "input": "width",
      "input_type": "keyvalue",
      "output": "package_width",
      "output_type": "keyvalue",
      "default_value": "",
      "fields": [],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "field_package_weight",
      "scope": "global",
      "input": "weight",
      "input_type": "keyvalue",
      "output": "package_weight",
      "output_type": "keyvalue",
      "default_value": false,
      "fields": [],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "obj_package_dimensions",
      "scope": "local",
      "input": "dimensions",
      "input_type": "object",
      "output": "package_dimensions",
      "output_type": "object",
      "default_value": "",
      "fields": ["field_package_height", "field_package_length","field_package_width", "field_package_weight"],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "arr_image_options",
      "scope": "local",
      "input": "images",
      "input_type": "arr",
      "output": "image_options",
      "output_type": "arr",
      "default_value": "",
      "fields": ["obj_image_option"],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "obj_image_option",
      "scope": "local",
      "input": "",
      "input_type": "obj",
      "output": "",
      "output_type": "obj",
      "default_value": "",
      "fields": ["field_image_link"],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "field_image_link",
      "scope": "local",
      "input": "src",
      "input_type": "keyvalue",
      "output": "link",
      "output_type": "keyvalue",
      "default_value": "",
      "fields": [],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "obj_description",
      "scope": "local",
      "input": "description",
      "input_type": "keyvalue",
      "output": "description",
      "output_type": "obj",
      "default_value": "",
      "fields": ["obj_description_data"],
      "execution_type": "nested",
      "helper_function": []
    },
    {
      "id": "obj_description_data",
      "scope": "local",
      "input": "description",
      "input_type": "keyvalue",
      "output": "data",
      "output_type": "keyvalue",
      "default_value": "",
      "fields": [],
      "execution_type": "nested",
      "helper_function": []
    }
  ]
}

```
##### Available Sample 1 fields are:

- **name** (*str*, required)
    - The human-readable name of the mapper
- **description** (*str*)
    - Extended description for the mapper, in reStructuredText
- **fields** (*obj*)
    - A field mapping describes how a persistent field maps to the database
    - The following table shows the sample mapper fields

- **field_sku_code**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | field_sku_code |
| scope | *str* | local |
| input | *str* | sku |
| input_type | *str* | keyvalue |
| output | *str* | sku_code |
| output_type | *str* | keyvalue |
| default_value | *str* | "sku-852052" |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **field_product_name**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | field_product_name |
| scope | *str* | local |
| input | *str* | name |
| input_type | *str* | keyvalue |
| output | *str* | product_name |
| output_type | *str* | keyvalue |
| default_value | *str* | "shirt" |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **field_package_height**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | field_package_height |
| scope | *str* | local |
| input | *str* | height |
| input_type | *str* | keyvalue |
| output | *str* | package_height |
| output_type | *str* | keyvalue |
| default_value | *str* | "5.6" |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **field_package_length**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | field_package_length |
| scope | *str* | local |
| input | *str* | length |
| input_type | *str* | keyvalue |
| output | *str* | package_length |
| output_type | *str* | keyvalue |
| default_value | *str* | "10" |
| fields  | *list(obj)* | [{"fields":"values"}] |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **field_package_width**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | field_package_width |
| scope | *str* | local |
| input | *str* | width |
| input_type | *str* | keyvalue |
| output | *str* | package_width |
| output_type | *str* | keyvalue |
| default_value | *str* | "10" |
| fields  | *list(obj)* | [{"fields":"values"}] |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **field_package_weight**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | field_package_weight |
| scope | *str* | local |
| input | *str* | weight |
| input_type | *str* | keyvalue |
| output | *str* | package_weight |
| output_type | *str* | keyvalue |
| default_value | *str* | "10" |
| fields  | *list(obj)* | [{"fields":"values"}] |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **obj_package_dimensions**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | obj_package_dimensions |
| scope | *str* | local |
| input | *str* | dimensions |
| input_type | *str* | object |
| output | *str* | package_dimensions |
| output_type | *str* | object |
| default_value | *str* | "10" |
| fields  | *list(obj)* | field_package_height", "field_package_length","field_package_width", "field_package_weight |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **arr_image_options**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | arr_image_options |
| scope | *str* | local |
| input | *str* | images |
| input_type | *str* | arr |
| output | *str* | image_options |
| output_type | *str* | arr |
| default_value | *str* | "10" |
| fields  | *list(obj)* | obj_image_option |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **obj_image_option**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | obj_image_option |
| scope | *str* | local |
| input | *str* | images |
| input_type | *str* | obj |
| output | *str* | image_options |
| output_type | *str* | obj |
| default_value | *str* | "asdfghjk.jpg" |
| fields  | *list(obj)* | field_image_link |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **field_image_link**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | field_image_link |
| scope | *str* | local |
| input | *str* | src |
| input_type | *str* | keyvalue |
| output | *str* | link |
| output_type | *str* | keyvalue |
| default_value | *str* | "asdfghjk.jpg" |
| fields  | *list(obj)* | field_image_link |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **obj_description**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | obj_description |
| scope | *str* | local |
| input | *str* | description |
| input_type | *str* | keyvalue |
| output | *str* | description |
| output_type | *str* | obj |
| default_value | *str* | "" |
| fields  | *list(obj)* | obj_description_data |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |

- **obj_description_data**

 | fields | Datatype | Sample Values |
| ------------- | ------------- | ------------- |
| id | *str* | obj_description_data |
| scope | *str* | local |
| input | *str* | description |
| input_type | *str* | keyvalue |
| output | *str* | data |
| output_type | *str* | keyvalue |
| default_value | *str* | "" |
| fields  | *list(obj)* | obj_description_data |
| execution_type | *str* | "nested" |
| helper_function | *list(obj)* | [] |








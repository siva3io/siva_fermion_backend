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

# Manifest File

The manifest file serves to declare a go package as separate module/app and to specify app metadata.

It is a file called manifest.json and contains a single json dictionary, where each key specifies app metadatum.

## Sample

```json
{
  "name": "An App",
  "version": "1.0.1.0.0",
  "author": "Author Name",
  "support": "sales@eunimart.com",
  "category": "Category",
  "summary": "Summary of the App",
  "description": "Description text",
  "website": "https://eunimart.com/",
  "license": "OPL-1",
  "data": ["data files always loaded at installation"],
  "demo": [
    "data files containing optionally loaded demonstration data"
  ],
  "currency": "EUR",
  "price": 279.0,
  "uninstall_hook": "uninstall_hook",
  "application": true
}
```

## Available manifest fields are:

- **name** (*str*, required)
    - The human-readable name of the app
- **version** (*str*)
    - This app’s version, should follow [semantic versioning](https://semver.org/) rules
- **author** (*str*)
  - Name of the app author
- **support** (*str*)
  - Email of the support team/author
- **category** (*str*, default: Uncategorized)
  - Classification category within Platform, rough business domain for the app.
  - Although using [existing categories]() is recommended, the field is freeform and unknown categories are created on-the-fly. Category hierarchies can be created using the separator / e.g. Foo / Bar will create a category Foo, a category Bar as child category of Foo, and will set Bar as the app’s category.
    When a app is installed, all of its dependencies are installed before it. Likewise dependencies are loaded before a app is loaded.
- **summary** (*str*)
  - A small summary about the app.
- **description** (*str*)
    - Extended description for the app, in reStructuredText
- **website** (*str*)
  - Website URL for the app author
- **license** (*str*, default: LGPL-3)
  - Distribution license for the app. Possible values:
    - GPL-2
    - GPL-2 or any later version
    - GPL-3
    - GPL-3 or any later version
    - AGPL-3
    - LGPL-3
    - Other OSI approved licence
    - EEEL-1 (Eunimart Enterprise Edition License v1.0)
    - EPL-1 (Eunimart Proprietary License v1.0)
    - Other proprietary
- **data** (*list(str)*)
  - List of data files which must always be installed or updated with the app. A list of paths from the app root directory
- **demo** (*list(str)*)
  - List of data files which are only installed or updated in demonstration mode
- **application** (*bool*, default: true)
  - Whether the app should be considered as a fully-fledged application (true) or is just a technical app (false) that provides some extra functionality to an existing application app.
- **external_dependencies** (*dict(key=list(str))*, required)
  - A dictionary containing go and/or binary dependencies.
- **maintainer** (*str*)
  - Person or entity in charge of the maintenance of this module, by default it is assumed that the author is the maintainer.
- **active** (*bool*, default: true)
  - The app is disabled or not
- **currency** (*str*)
  - The Currency code as per [ISO 4217](https://en.wikipedia.org/wiki/ISO_4217)
- **price** (*float*, default= 0.0)
  - The price of app in case if its paid.
- **** (**, required)
  - 
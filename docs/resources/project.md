---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wistia_project Resource - terraform-provider-wistia"
subcategory: ""
description: |-
  A Wistia project. See the API documentation https://wistia.com/support/developers/data-api#projects for more details.
---

# wistia_project (Resource)

A Wistia project. See the [API documentation](https://wistia.com/support/developers/data-api#projects) for more details.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The project's display name.

### Optional

- **anonymous_can_download** (Boolean) A boolean indicating whether or not anonymous downloads are enabled for the project.
- **anonymous_can_upload** (Boolean) A boolean indicating whether or not anonymous uploads are enabled for the project.
- **id** (String) The ID of this resource.
- **public** (Boolean) A boolean indicating whether the project is available for public (anonymous) viewing.

### Read-Only

- **created** (String) The date that the project was originally created.
- **description** (String) The project's description.
- **hashed_id** (String) A private hashed ID, uniquely identifying the project within the system.
- **media_count** (Number) The number of different medias that have been uploaded to the project.
- **public_id** (String) If the project is public, this field contains a string representing the ID used for referencing the project in public URLs.
- **updated** (String) TThe date that the project was last updated.



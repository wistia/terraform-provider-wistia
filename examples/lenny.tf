provider "wistia" {
  environment = "staging"
}

resource "wistia_project" "my_first_terraformed_project" {
  name = "Lenny's Terraformed Wistia Project"
  anonymous_can_download = true
  anonymous_can_upload = false
  public = false
}

resource "wistia_media" "my_first_terraformed_media" {
  name = "Lenny terraforms Mars"
  url = "https://embed-ssl.wistia.com/deliveries/57be37488565a5e51351e5ecebcea610.bin"
  project_id = wistia_project.my_first_terraformed_project.hashed_id
}

resource "wistia_media" "my_second_terraformed_media" {
  name = "Lenny terraforms Mars again"
  file = "lenny-mars.mp4"
  project_id = wistia_project.my_first_terraformed_project.hashed_id
}

resource "wistia_project" "my_second_terraformed_project" {
  name = "Lenny's Second Terraformed Wistia Project"
  anonymous_can_download = false
  anonymous_can_upload = false
  public = false
}

resource "wistia_media" "my_third_terraformed_media" {
  name = "Lenny terraforms Wistia HQ"
  file = "lenny-wistia-hq.mp4"
  project_id = wistia_project.my_second_terraformed_project.hashed_id
}
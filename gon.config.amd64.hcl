source = [
  "dist/classeviva-macos-amd64_darwin_amd64_v1/classeviva"
]
bundle_id = "dev.zmoog.classeviva"

sign {
  application_identity = "Developer ID Application: Maurizio Branca (47UB39QJQM)"
}

# Ask Gon for zip output to force notarization process to take place.
# The CI will ignore the zip output, using the signed binary only.
zip {
  output_path = "unused.amd64.zip"
}
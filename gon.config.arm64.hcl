source = [
  "dist/classeviva-macos-arm64_darwin_arm64/classeviva"
]
bundle_id = "dev.zmoog.classeviva"

sign {
  application_identity = "Developer ID Application: Maurizio Branca (47UB39QJQM)"
}

# Ask Gon for zip output to force notarization process to take place.
# The CI will ignore the zip output, using the signed binary only.
zip {
  output_path = "unused.arm64.zip"
}
# OpenAPI Specification

This directory contains the OpenAPI 3.0 specification for the Classeviva API.

## What is this?

The `openapi.yaml` file is a machine-readable API specification that describes all the endpoints, request/response formats, and data models for the Classeviva school portal API.

**Note**: This is an unofficial specification created by reverse engineering and community documentation. The official Classeviva API is not publicly documented by Spaggiari.

## What can I do with it?

### 1. View Interactive Documentation

You can use various tools to view the API documentation interactively:

#### Using Swagger UI Online
1. Go to https://editor.swagger.io/
2. Click "File" → "Import File"
3. Upload the `openapi.yaml` file
4. Browse the interactive documentation

#### Using Redoc
```bash
npx @redocly/cli preview-docs openapi.yaml
```

Then open http://localhost:8080 in your browser.

### 2. Generate API Clients

Generate a client library in your preferred programming language:

#### Python
```bash
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g python \
  -o ./generated/python-client
```

#### JavaScript/TypeScript
```bash
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g typescript-axios \
  -o ./generated/typescript-client
```

#### Java
```bash
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g java \
  -o ./generated/java-client
```

#### Other Languages
OpenAPI Generator supports [50+ languages](https://openapi-generator.tech/docs/generators). Replace the `-g` flag with your target language.

### 3. Validate API Requests

Use the specification to validate API requests and responses in your tests:

```bash
npm install -g @redocly/cli
redocly lint openapi.yaml
```

### 4. Mock API Server

Create a mock server for testing without hitting the real API:

```bash
npx @stoplight/prism-cli mock openapi.yaml
```

## API Overview

The Classeviva API provides access to:

- **Authentication**: Login and session management
- **Grades**: Student grades and evaluations
- **Agenda**: Homework, events, and assignments
- **Noticeboard**: School announcements and circulars
- **Absences**: Student absences and late arrivals
- **Lessons**: Lesson records and topics
- **Calendar**: School calendar
- **Student Profile**: Personal information
- **Subjects**: Course subjects and teachers
- **Academic Periods**: Trimesters, semesters
- **Notes**: Teacher notes and warnings
- **Didactics**: Teaching materials
- **Schoolbooks**: Required textbooks

## Authentication

All API requests require these headers:

```
Z-Dev-Apikey: Tg1NWEwNGIgIC0K
User-Agent: CVVS/std/4.2.3 Android/12
Content-Type: application/json
```

For authenticated endpoints, you also need:

```
Z-Auth-Token: <token from login response>
```

### Login Flow

1. Call `POST /auth/login` with credentials:
   ```json
   {
     "uid": "username",
     "pass": "password",
     "ident": null
   }
   ```

2. Extract the `token` and student `ident` from the response

3. Extract the numeric student ID from `ident` (e.g., "G9123456R" → "9123456")

4. Use the token in subsequent requests with the `Z-Auth-Token` header

5. Use the student ID in endpoint paths (e.g., `/students/{studentId}/grades`)

## Contributing

If you find any issues or want to add missing endpoints, please:

1. Check the unofficial documentation at [Classeviva-Official-Endpoints](https://github.com/Lioydiano/Classeviva-Official-Endpoints)
2. Open an issue or pull request

## Resources

- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3)
- [OpenAPI Generator](https://openapi-generator.tech/)
- [Swagger Editor](https://editor.swagger.io/)
- [Redocly CLI](https://redocly.com/docs/cli/)
- [Unofficial API Documentation](https://github.com/Lioydiano/Classeviva-Official-Endpoints)

---
name: analyze-code-flow
description: Analyzes and explains code flow, execution paths, function calls, and system architecture. Use when the user asks to understand how code works, trace execution flow, explain functions, analyze system behavior, interpret logic, or understand code/system architecture.
---

# Analyze Code Flow

This skill helps systematically analyze and explain how code executes, including function calls, data flow, control flow, and system interactions.

## When to Use This Skill

Use this skill when the user asks about:
- "How does X work?"
- "Trace the flow of..."
- "Explain this function/component"
- "What happens when..."
- "Walk me through the execution"
- Understanding system architecture or module interactions

## Analysis Approach

Provide **two-level analysis** by default:

### Level 1: Overview (Always provide first)
High-level understanding of the main flow:
1. Entry point and trigger
2. Key steps in sequence
3. Major components/functions involved
4. Final outcome or return value

### Level 2: Detailed Analysis (Provide after overview)
Deep dive into the implementation:
1. Line-by-line execution logic
2. Variable transformations
3. Conditional branches and edge cases
4. Dependencies and side effects
5. Error handling
6. Performance considerations

## Analysis Structure

Use this template:

```markdown
## Flow Overview

**Entry Point**: [Function/endpoint/trigger]
**Purpose**: [What this code accomplishes]

**Main Flow**:
1. [Step 1 - high level]
2. [Step 2 - high level]
3. [Step 3 - high level]

**Key Components**:
- `ComponentName`: [Role in the flow]
- `FunctionName`: [What it does]

**Outcome**: [What happens at the end]

---

## Detailed Analysis

### Step 1: [Description]
**Location**: `filename:lines`
**What happens**:
- [Detailed explanation]
- [Variable changes]
- [Logic decisions]

**Code**:
[Show relevant code snippet]

### Step 2: [Description]
[Continue pattern...]

---

## Key Insights

- **Data Flow**: [How data transforms through the process]
- **Control Flow**: [Conditional paths and branches]
- **Side Effects**: [External interactions, state changes]
- **Edge Cases**: [Special conditions handled]
```

## Best Practices for Analysis

### 1. Start Broad, Then Go Deep
- Begin with the big picture
- Then zoom into implementation details
- Don't overwhelm with details upfront

### 2. Follow Execution Order
- Trace code in the order it executes
- Number steps sequentially
- Show how control flows between functions

### 3. Highlight Key Decision Points
- Identify if/else branches
- Explain loop conditions
- Show how errors are handled

### 4. Track Data Transformations
- Show input values
- Explain how data changes
- Highlight output or side effects

### 5. Reference Code Precisely
Use code references with line numbers:
```
startLine:endLine:filepath
```

### 6. Use Clear Terminology
- **Entry point**: Where execution begins
- **Call chain**: Sequence of function calls
- **Control flow**: How execution moves through code
- **Data flow**: How data is passed and transformed
- **Side effect**: External state changes (DB, files, API calls)

## Analysis Patterns

### Pattern 1: Function Analysis

```markdown
## Function: `functionName`

**Purpose**: [What it does]
**Parameters**: [List with types and purpose]
**Returns**: [Type and meaning]

**Flow**:
1. [Validate inputs]
2. [Process data]
3. [Handle errors]
4. [Return result]

**Dependencies**: [Functions/modules it calls]
**Side Effects**: [Any external changes]
```

### Pattern 2: API Endpoint Analysis

```markdown
## Endpoint: `METHOD /path`

**Purpose**: [What this endpoint does]
**Authentication**: [Requirements]
**Request**: [Body/params structure]
**Response**: [Expected output]

**Execution Flow**:
1. **Middleware**: [Auth, validation, etc.]
2. **Handler**: [Main logic]
3. **Database**: [Queries executed]
4. **Response**: [What gets returned]

**Error Scenarios**: [How failures are handled]
```

### Pattern 3: System Flow Analysis

```markdown
## System Flow: [Feature Name]

**Trigger**: [What initiates this flow]

**Component Interaction**:
1. **Component A** → [Action] → **Component B**
2. **Component B** → [Action] → **Component C**
3. [Continue chain...]

**State Changes**: [What gets modified]
**External Systems**: [APIs, databases, services involved]
```

### Pattern 4: Conditional Flow Analysis

```markdown
## Conditional Logic: [Description]

**Conditions Evaluated**:

**Branch 1**: If [condition]
- [What happens]
- [Code path taken]

**Branch 2**: Else if [condition]
- [What happens]
- [Code path taken]

**Branch 3**: Else
- [Default behavior]

**Critical Logic**: [Important decisions or edge cases]
```

## Questions to Answer During Analysis

As you analyze code, address these questions:

1. **Entry & Exit**:
   - Where does execution start?
   - How does it end or what does it return?

2. **Purpose**:
   - What problem does this solve?
   - Why does it exist?

3. **Data**:
   - What inputs does it receive?
   - How is data transformed?
   - What outputs are produced?

4. **Control**:
   - What are the main execution paths?
   - What conditions cause different behaviors?
   - Are there loops or recursion?

5. **Dependencies**:
   - What other code does this call?
   - What external systems does it interact with?

6. **Error Handling**:
   - What can go wrong?
   - How are errors handled?

7. **Side Effects**:
   - Does it modify state?
   - Does it make external calls?

## Tools for Analysis

### Read Code
Use `Read` tool to examine source files:
- Start with the entry point
- Follow function calls
- Read dependencies as needed

### Search for References
Use `Grep` to find:
- Where functions are called
- How types/interfaces are used
- Related code patterns

### Semantic Search
Use `SemanticSearch` for:
- Finding related functionality
- Understanding system architecture
- Discovering implementation patterns

## Common Analysis Scenarios

### Scenario 1: "How does login work?"
1. Find the login endpoint/function
2. Trace authentication flow
3. Identify credential validation
4. Show session/token creation
5. Explain error handling

### Scenario 2: "What happens when a user submits a form?"
1. Identify form handler
2. Show validation logic
3. Trace data processing
4. Follow database operations
5. Explain response generation

### Scenario 3: "Explain this error"
1. Find where error originates
2. Trace backwards to root cause
3. Identify conditions that trigger it
4. Show error propagation
5. Suggest fixes

### Scenario 4: "How do these modules interact?"
1. Map component boundaries
2. Show data exchange
3. Identify coupling points
4. Explain communication patterns
5. Highlight dependencies

## Output Quality Guidelines

### Be Precise
- ✅ "The `validateUser` function checks if email exists in database"
- ❌ "It validates the user somehow"

### Show, Don't Just Tell
- ✅ Include relevant code snippets with line numbers
- ❌ Only describe without showing code

### Explain the "Why"
- ✅ "This checks for nil to prevent panic when accessing properties"
- ❌ "This checks for nil"

### Connect the Pieces
- ✅ "After validation succeeds, the token is passed to `createSession` which stores it in Redis"
- ❌ List functions without showing relationships

### Highlight Important Details
- Point out clever solutions
- Flag potential issues
- Note performance implications
- Identify security considerations

## Progressive Disclosure

For large or complex systems:

1. **Start Simple**: Give the 30-second overview
2. **Build Up**: Add more detail layer by layer
3. **Offer Depth**: "Would you like me to explain [specific part] in more detail?"
4. **Stay Focused**: Don't analyze everything at once

## Anti-Patterns to Avoid

❌ **Reading code out loud**: "Line 5 declares a variable called x"
✅ **Explain purpose**: "The function accumulates user scores in the `total` variable"

❌ **Overwhelming details upfront**: Starting with line-by-line analysis
✅ **Progressive detail**: Overview first, then deeper dive

❌ **Ignoring context**: Analyzing in isolation
✅ **Show connections**: How it fits in the larger system

❌ **Assuming knowledge**: Using unexplained jargon
✅ **Clear explanations**: Define terms and concepts

❌ **Missing the forest for trees**: Getting lost in implementation details
✅ **Balance**: Both high-level flow and important details

## Example Analysis Output

Here's what a good analysis looks like:

```markdown
## Flow Overview: User Registration

**Entry Point**: `POST /api/auth/register` → `authController.register()`
**Purpose**: Create new user account with email verification

**Main Flow**:
1. Validate registration data (email, password)
2. Check if user already exists
3. Hash password with bcrypt
4. Create user record in database
5. Generate verification token
6. Send verification email
7. Return success response

**Key Components**:
- `authController`: Handles HTTP request/response
- `userService`: Business logic for user operations
- `emailService`: Sends verification email
- `database`: PostgreSQL user table

**Outcome**: User created with pending verification status

---

## Detailed Analysis

### Step 1: Request Validation
**Location**: `auth/controller.go:45-52`

**What happens**:
- Extracts email and password from request body
- Validates email format using regex
- Checks password meets complexity requirements (min 8 chars, etc.)
- Returns 400 error if validation fails

**Code**:
```45:52:auth/controller.go
func (c *Controller) register(ctx *gin.Context) {
    var req RegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        return ctx.JSON(400, ErrorResponse{Message: "Invalid request"})
    }
    // Validation continues...
}
```

[Continue with remaining steps...]

---

## Key Insights

- **Data Flow**: Plain password → bcrypt hash → database
- **Security**: Passwords never stored in plain text, verification token expires
- **Error Cases**: Duplicate email returns 409, invalid data returns 400
- **Side Effects**: Database write, email sent to external service
```

## Summary

When analyzing code flow:
1. Always start with overview, then offer details
2. Follow execution order chronologically
3. Explain both what happens and why
4. Show relevant code with precise references
5. Track data transformations and side effects
6. Identify key decision points and edge cases
7. Connect components and explain interactions
8. Use clear, consistent terminology throughout

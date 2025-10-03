# Buildathon Fall 2025 - Project Research

## Project Overview
- **Event Type**: Company-wide buildathon
- **Duration**: 3 Fridays
- **Format**: Teams work on small demoable projects
- **Status**: Research Phase

## Project Information
**Project Goal**: Build an MCP (Model Context Protocol) tool that can call Dispatch's order creation and estimate mutations

**Project Type**: POC/MVP demonstration to show how easy it is to integrate MCP tools with Dispatch systems

**Project Description** (from Chris's project tracking):
Customers struggle with complex delivery options and pricing. Our MCP enables conversational order placement - "I need this from A to B by 5pm" - and intelligently presents options and places orders for them.

**Project Impact** (from Chris's project tracking):
Transforms customer experience by replacing complex shipping forms and decision trees with natural conversation. Customers can describe delivery needs in plain English and receive recommendations with transparent pricing. Streamlines internal order processing by automatically capturing requirements, validating delivery constraints, and routing orders through existing systems.

**Key Objectives**:
- Demonstrate MCP tool integration with Dispatch APIs
- Focus on order creation and estimate mutations
- Keep it simple and POC-ish for buildathon timeframe
- Make each step clear and repeatable
- Enable conversational order placement with natural language
- Transform complex delivery forms into simple conversations

## Team Information
**Team Members**:
- **Chris Hayes** - Engineering Manager (Project Lead)
- **Camron Wood** - Backend/Go Engineer (Technical Lead for Go service)
- **Julia** - Frontend Developer (UI development)
- **Tyler** - Pre-sales Engineer (Dispatch historian, domain expertise)

**Project Tracking**:
- **Team**: Chris, Camron, Tyler, Julia
- **Tools**: Claude Code
- **Project Tracking Sheet**: https://docs.google.com/spreadsheets/d/1MJomCDNBCV1xQjS3aBt1yBgsqZE629MqUhyG7bTjf7Y/edit?pli=1&gid=0#gid=0

**Roles & Responsibilities**:
- Camron: Lead on Go service development
- Julia: UI development and frontend
- Tyler: Domain knowledge, explaining business logic (cargo vans vs trucks, etc.)
- Chris: Project oversight and guidance

## Technical Requirements
**Core Technology Stack**:
- **Language**: Go
- **Framework**: MCP (Model Context Protocol)
- **Integration**: Dispatch order creation and estimate mutations APIs

**Technical References**:
- Chris's existing MCP server for Jira APIs: https://github.com/dispatchitinc/go-jira/blob/feature/mcp/internal/mcp/server.go
- MCP Go framework: https://github.com/mark3labs/mcp-go

**Target GraphQL Mutations** (Platform & Integrations team):
- **Create Estimate**: https://github.com/dispatchitinc/ms-monolith/blob/main/app/graphql/mutations/order/create_estimate.rb
- **Create Order**: https://github.com/dispatchitinc/ms-monolith/blob/main/app/graphql/mutations/order/create_order.rb

**Authentication Strategy**:
- MVP approach: Embed auth tokens in environment variables (hacky but simple for POC)
- Future consideration: Proper integration with Dispatch auth model

## Research Notes
**Key Technical Challenges**:
1. **Authentication Integration**: Hardest part will be hooking into Dispatch's auth model
2. **API Integration**: Understanding Dispatch's order creation and estimate mutation endpoints
3. **MCP Implementation**: Following Chris's Jira MCP server as a reference

**Next Steps**:
- Sync with Chris before Friday to clarify project approach
- Reach out to Jon Reyes for pointers and tips on tools
- Review Chris's existing MCP server implementation
- Understand Dispatch's order/estimate API structure

**Pre-Buildathon Homework** (from Chris):
1. **Familiarize with MCP-go structure**: Review https://github.com/mark3labs/mcp-go
2. **Study target mutations**: 
   - Create Estimate: https://github.com/dispatchitinc/ms-monolith/blob/main/app/graphql/mutations/order/create_estimate.rb
   - Create Order: https://github.com/dispatchitinc/ms-monolith/blob/main/app/graphql/mutations/order/create_order.rb
3. **Goal**: Understand what these mutations do for Friday head start

## MCP Implementation Pattern Analysis

**Based on research of existing MCP implementations:**

**Common MCP Server Structure**:
- **MCP Server**: Central orchestrator receiving requests from AI clients
- **Capabilities/Tools**: Specialized adapters for external APIs (Jira, Dispatch, etc.)
- **Client Interface**: Natural language processing for AI interactions
- **Authentication**: Environment-based configuration for API credentials

**Key MCP Components**:
1. **Server Setup**: Using mcp-go framework for Go-based MCP servers
2. **Tool Definition**: Functions that perform specific actions (create_estimate, create_order)
3. **API Integration**: HTTP/GraphQL clients for external service communication
4. **Error Handling**: Comprehensive error management and logging
5. **Configuration**: Environment variables for credentials and settings

**MCP Server Pattern**:
```go
// Typical MCP server structure
type MCPServer struct {
    tools    map[string]Tool
    client   APIClient
    config   Config
}

type Tool interface {
    Execute(params map[string]interface{}) (interface{}, error)
}
```

**Authentication Strategy**:
- Environment variables for API tokens (MVP approach)
- Secure credential management
- Configurable endpoints and settings

## Resources & References
**Code References**:
- Chris's Jira MCP Server: https://github.com/dispatchitinc/go-jira/blob/feature/mcp/internal/mcp/server.go
- MCP Go Framework: https://github.com/mark3labs/mcp-go

**Team Contacts**:
- **Jon Reyes**: For pointers and tips on tools (mentioned by Camron)
- **Chris Hayes**: Project lead and MCP experience
- **Tyler**: Domain expertise on Dispatch business logic

## Timeline & Milestones
**Current Status**: Research and planning phase
**Upcoming**:
- Sync with Chris before Friday to clarify project approach
- Reach out to Jon Reyes for tool guidance
- Review existing MCP implementations
- Understand Dispatch API structure

**Buildathon Schedule**: 3 Fridays (specific dates TBD)

## Meeting Notes: Camron & Chris Discussion

### Key Discussion Points

**Challenges with Documentation and Standardization**:
- Camron expressed difficulty finding clear documentation and linear paths for operations (deployment, sub-environment processes)
- Chris agreed, critiquing Go for not being opinionated, leading to various paradigms and lack of standardization in microservices
- Both agreed on need for frameworks and documentation to establish best practices and consistent patterns
- Go's readability sometimes leads developers to prioritize speed over proper structuring

**Buildathon Project Vision**:
- **Goal**: Create a conversational interface for placing dispatch orders and providing pricing estimates
- **Interface Type**: "Smart form" or "automated order form" where LLM guides user through collecting necessary information
- **Priority**: Proof of concept and MVP - functionality over perfection
- **Value**: Similar to how they personally use LLMs for Jira ticket management or understanding codebases

**Technical Approach Decisions**:
- **Focus**: Build MCP server in Go that calls external APIs for order creation and cost estimates
- **Avoid**: GraphQL mutations might be too complex for project scope
- **Authentication**: Consider basic HTTP authentication if OAuth 2 proves too time-consuming
- **Architecture**: MCP server as the core, with external API integration

**Team Roles & Responsibilities**:
- **Camron**: Go development for MCP server
- **Julia**: Design LLM prompt incorporating "dispatchisms" and industry-specific knowledge (vehicle capacities, dimensions)
- **Julia**: Build interface (ideally web-based with dispatch logo, but Claude desktop acceptable for POC)
- **Tyler**: Domain expertise and business logic

**Strategic Objectives**:
- Prove that such projects are not difficult (contrast to perceived slowness in Greg's team)
- Deliver tangible proof of concept
- Establish core mechanics: order placement, pricing, response parsing
- UI is "cherry on top" - focus on backend functionality first

**Next Steps & Action Items**:
1. **Camron**: Schedule Friday morning team meeting for alignment
2. **Camron**: Start by analyzing necessary API endpoints and researching authentication methods
3. **Chris**: Post any information or resources related to the project in the channel
4. **Julia**: Work on writing the prompt with dispatchisms
5. **Julia**: Build an interface (web-based preferred, but simple approach acceptable)
6. **Team**: Evaluate both APIs first to determine easiest path for authentication

**Pivot Strategy**:
- Quick pivoting if roadblocks encountered (e.g., authentication issues)
- Focus on core functionality over perfect implementation
- MVP approach with environment variable auth tokens (hacky but simple for POC)

## Questions & Next Steps
**Immediate Action Items**:
1. Schedule sync with Chris before Friday
2. Contact Jon Reyes for tool guidance
3. Review Chris's Jira MCP server implementation
4. Research Dispatch's order creation and estimate mutation APIs
5. Understand authentication requirements

**Pre-Buildathon Prep** (Chris's homework):
1. Study MCP-go framework structure
2. Review the two target GraphQL mutations (create_estimate.rb and create_order.rb)
3. Understand what these mutations do and their parameters

**Open Questions**:
- What are the specific Dispatch API endpoints for orders and estimates?
- How does Dispatch's current auth model work?
- What's the exact structure of order creation and estimate mutations?
- What UI components does Julia need to build?
- What domain knowledge does Tyler need to provide?

---
*Last Updated: [Date will be updated as information is added]*

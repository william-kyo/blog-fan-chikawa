## Project Planning File Management
1. Use `.claude/session.md` for temporary task planning during current session
2. Long-term planning is maintained in `TODO.md`
3. Project documentation goes in `docs/` directory
4. Claude working files are not committed to repository

## Standard Workflow
1. First check `TODO.md` to understand overall project planning
2. Think through the problem, read the codebase for relevant files, and write a plan to `.claude/session.md`
3. The plan should have a list of todo items that you can check off as you complete them
4. Before you begin working, check in with me and I will verify the plan
5. Then, begin working on the todo items, marking them as complete as you go
6. Please every step of the way just give me a high level explanation of what changes you made
7. Make every task and code change you do as simple as possible. We want to avoid making any massive or complex changes. Every change should impact as little code as possible. Everything is about simplicity
8. Finally, update `TODO.md` with completed long-term tasks and add a review section to `.claude/session.md`

## Git Commit Guidelines
- NEVER include Claude Code attribution or Co-Authored-By: Claude information in commit messages
- Use conventional commit format: `<type>: <description>`
- Keep commit messages concise and focused on the actual changes
- Author should always be "kyo <243075803@qq.com>"

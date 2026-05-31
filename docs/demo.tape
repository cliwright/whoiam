Output docs/demo.gif

Require echo

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 600
Set Theme "Catppuccin Mocha"
Set Padding 20
Set TypingSpeed 60ms

# Start a session
Type "aws-vault exec whoiam.demo.dev"
Enter
Sleep 5s

# Show current identity
Type "whoiam"
Enter
Sleep 2s

# Add account mappings to config
Sleep 1s
Type "cat .whoiam/whoiam.yaml"
Enter
Sleep 2s

# Set expected environment
Sleep 1s
Type "whoiam set whoiam.demo.dev --global"
Enter
Sleep 1s

# Check status
Type "whoiam status"
Enter
Sleep 2s

# Validate against expected account
Sleep 1s
Type "whoiam validate"
Enter
Sleep 2s
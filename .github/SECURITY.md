# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security issues seriously. We appreciate your efforts to responsibly disclose your findings.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **GitHub Security Advisories**: Use [GitHub's private vulnerability reporting](https://github.com/ageha734/mmdc/security/advisories/new) feature.

2. **Email**: Send an email to the maintainers (check the repository for contact information).

### What to Include

Please include the following information in your report:

- Type of issue (e.g., buffer overflow, command injection, privilege escalation)
- Full paths of source file(s) related to the issue
- Location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: 7 days
  - High: 14 days
  - Medium: 30 days
  - Low: 90 days

### What to Expect

1. **Acknowledgment**: We will acknowledge receipt of your vulnerability report.
2. **Communication**: We will keep you informed about the progress toward a fix.
3. **Credit**: We will credit you in the release notes (unless you prefer to remain anonymous).
4. **Disclosure**: We will coordinate with you on the public disclosure timeline.

## Security Best Practices

When using mmdc, please follow these best practices:

### Input Validation

- Always validate Mermaid diagram input before processing
- Be cautious when processing untrusted input files
- Use the `--quiet` flag in automated pipelines to avoid exposing sensitive information in logs

### File System Security

- Ensure output directories have appropriate permissions
- Be careful when using stdin input in automated scripts
- Avoid running mmdc with elevated privileges when possible

### Browser Security

mmdc uses chromedp to render diagrams in a headless Chrome instance. The following security measures are in place:

- Temporary user data directories are created and cleaned up for each session
- Headless mode is always enabled
- No persistent browser data is stored

### Dependency Security

We use Dependabot to monitor and update dependencies for known vulnerabilities. We recommend:

- Keeping mmdc updated to the latest version
- Checking the [Security Advisories](https://github.com/ageha734/mmdc/security/advisories) page regularly

## Security Features

### Code Signing

Release binaries are signed and checksums are provided. Always verify downloads:

```bash
# Verify checksum
sha256sum -c checksums.txt

# Or manually
sha256sum mmdc-linux-amd64.tar.gz
```

### SBOM

Software Bill of Materials (SBOM) is generated for each release and available in the release assets.

## Known Security Considerations

1. **Browser Execution**: mmdc executes a headless Chrome browser. While sandboxed, this carries inherent risks when processing untrusted input.

2. **File Access**: The tool reads input files and writes output files. Ensure appropriate file permissions.

3. **External Resources**: Mermaid diagrams may reference external resources (icons, fonts). These are loaded by the browser during rendering.

## Acknowledgments

We thank all security researchers who have responsibly disclosed vulnerabilities to us.

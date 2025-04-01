# Configuration

## Logger

A fixed *production* configuration was used for the logger,
so debug messages were not displayed at all.


## Environment Variables

Unlike as documented in the *README.md* file,
all environment variables must be prefixed with **APP_**
*(as instructed via envconfig.Process)*.

So instead of **APP_LOG_LEVEL** as documented, it actually must be injected as:
**APP_APP_LOG_LEVEL**

# API Specs

- Wrong URI path in *get winning ads* example response

# Models

Models haven been place in one file, rather than domain specific files.
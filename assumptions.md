# Winning ads

- Ads with a high bid get prioritized as they generate more revenue for the company.
- Assuming the fictional categories in order of highest to lowest score: electronics, sale, fashion, travel, finance
- Ads get an additional score based on the category, so that ads place with the category electronics (see above) get boosted.
- Given a successful conversion (any type) was executed, the budget gets decreased by the related bid.

# Tracking

- The task requirements only mention to record user interactions, but no logic on how to use the tracking events later
  on.
  In this case I went with a simple tracking, given the opportunity to filter based on criteria such as: line_item_id,
  event type, etc.
- If the whole budget has been spent, the line item status gets updated to: **completed**
  Even in case the budget gets negative (due to no endpoint validation - fairness for the customer), there is no further
  action necessary.

# Code quality

- Linting is based on [golangci-lint](https://golangci-lint.run/) with standard configuration using the latest version.
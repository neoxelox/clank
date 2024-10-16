#!/bin/bash
set -eo pipefail

COLOR="danger"
for arg in "$@"; do
  case $arg in
    --success)
      COLOR="good"
      shift
      ;;
    *)
      ;;
  esac
done

MESSAGE=$(cat <<-EOF
{
    "attachments": [{
        "color": "${COLOR}",
        "mrkdwn_in": ["text"],
        "text": "*${MONIT_EVENT}* for *\`${MONIT_SERVICE}\`* at \`${MONIT_HOST}\`\n\`\`\`${MONIT_DESCRIPTION}\`\`\`"
    }]
}
EOF
)

curl -s -X POST --data-urlencode "payload=$MESSAGE" $MONIT_SLACK_WEBHOOK

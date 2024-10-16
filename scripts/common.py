def ensure_envfile(context, envfile):
    with open(envfile, "r+") as file:
        content = file.read()

        content = content.replace("#:replace {$COMMIT}", context.commit())

        file.seek(0)
        file.write(content)
        file.truncate()

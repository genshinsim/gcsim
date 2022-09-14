cd "./db"
yarn build

wrangler pages publish build --project-name gcsimdb --branch main
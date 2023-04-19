import type {GatsbyConfig} from "gatsby"

const config: GatsbyConfig = {
    siteMetadata: {
        title: `Astro`,
        siteUrl: `https://www.yourdomain.tld`,
    },
    // More easily incorporate content into your pages through automatic TypeScript type generation and better GraphQL IntelliSense.
    // If you use VSCode you can also use the GraphQL plugin
    // Learn more at: https://gatsby.dev/graphql-typegen
    graphqlTypegen: false,
    plugins: [
        `gatsby-plugin-sass`,
        `gatsby-plugin-image`,
        `gatsby-plugin-sharp`,
        `gatsby-transformer-sharp`,
        {
            resolve: `gatsby-plugin-manifest`,
            options: {
                icon: `src/assets/images/favicons/favicon.png`,
                icons: [
                    {
                        src: `src/assets/images/favicons/android-chrome-192x192.png`,
                        sizes: `192x192`,
                        type: `image/png`,
                    },
                    {
                        src: `src/assets/images/favicons/android-chrome-512x512.png`,
                        sizes: `512x512`,
                        type: `image/png`,
                    },
                    {
                        src: `src/assets/images/favicons/apple-touch-icon.png`,
                        sizes: `180x180`,
                        type: `image/png`,
                    },
                    {
                        src: `src/assets/images/favicons/favicon-16x16.png`,
                        sizes: `16x16`,
                        type: `image/png`,
                    },
                    {
                        src: `src/assets/images/favicons/favicon-32x32.png`,
                        sizes: `32x32`,
                        type: `image/png`,
                    }
                ],
            }
        }
    ],
}

export default config

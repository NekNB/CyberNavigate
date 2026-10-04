import { MetadataRoute } from "next";

export default function robots(): MetadataRoute.Robots {
  return {
    rules: {
      userAgent: "*",
      allow: ["/", "/articles", "/simulator"],
      disallow: "/scenario",
    },
    sitemap: "https://кибер-навигатор.рф/sitemap.xml",
  };
}

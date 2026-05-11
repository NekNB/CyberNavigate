window.onload = function () {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    urls: [
      {
        name: "User Service",
        url: "/specs/docs/user/user.swagger.json", // Ссылка должна совпадать с тем, что в Go коде
      },

      // TODO: Добавить новые пути
      {
        name: "Article Service",
        url: "/specs/docs/article-service/article.swagger.json",
      },
    ],
    dom_id: "#swagger-ui",
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "StandaloneLayout",
  });

  //</editor-fold>
};

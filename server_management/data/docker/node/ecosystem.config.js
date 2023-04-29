module.exports = {
  apps: [
    {
      name: "node",
      script: "npm",
      args: "start",
      watch: true,
      log_date_format: "YYYY-MM-DD HH:mm Z",
      // 개발환경시 적용될 설정 지정
      env: {
        NODE_ENV: "production",
        PORT: "80",
      },
      // 배포환경시 적용될 설정 지정
      env_production: {
        NODE_ENV: "production",
        PORT: "80",
      },
    },
  ],
};

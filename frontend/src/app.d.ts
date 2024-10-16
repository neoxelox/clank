declare global {
  namespace App {
    interface MetaData {
      title: string;
      description?: string;
      image?: string;
    }

    interface PageData {
      meta: MetaData;
    }

    interface Platform {
      env: Env;
      cf: CfProperties;
      ctx: ExecutionContext;
    }
  }
}

export {};

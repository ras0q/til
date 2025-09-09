import { env } from "$env/dynamic/private";
import type { PageServerLoad } from "./$types.d.ts";

export const load: PageServerLoad = () => {
  return {
    pwd: env.PWD,
  };
};

import * as entities from "$lib/entities";
import { persisted } from "svelte-local-storage-store";
import { writable } from "svelte/store";

export const LOCAL_STORAGE_PREFIX: string = "_clank_";

export const cookies = persisted<string | undefined>(LOCAL_STORAGE_PREFIX + "cookies", undefined);

export const signInCodeID = persisted<string | undefined>(LOCAL_STORAGE_PREFIX + "sign_in_code_id", undefined);
export const signInCodeState = persisted<string | undefined>(LOCAL_STORAGE_PREFIX + "sign_in_code_state", undefined);

export const user = writable<entities.User>();
export const organization = writable<entities.Organization>();
export const users = writable<entities.User[]>();
export const invitations = writable<entities.Invitation[]>();
export const products = writable<entities.Product[]>();
export const product = writable<entities.Product>();

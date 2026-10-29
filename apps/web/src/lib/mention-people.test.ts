import test from "node:test";
import assert from "node:assert/strict";
import { collectMentionPeople, mentionCandidatesForChannel } from "./chat/people.ts";

const user = (id: string, handle: string) => ({ id, handle, display_name: handle });

test("keeps non-recent ordinary workspace members available for mention highlighting", () => {
  const members = Array.from({ length: 30 }, (_, index) =>
    user(`member-${index}`, `member-${index}`),
  );
  const target = members[29];

  const people = collectMentionPeople(null, [], members, undefined);

  assert.equal(people.length, 30);
  assert.equal(people.find((person) => person.id === target.id)?.handle, target.handle);
});

test("suggests only bots bound to the current channel (DHC)", () => {
  const people = [
    { ...user("u1", "alex") },
    { ...user("b1", "kunde-a-am"), bound_channel_ids: ["chn_a"] },
    { ...user("b2", "kunde-b-am"), bound_channel_ids: ["chn_b"] },
    { ...user("b3", "frei") },
  ];
  const handles = mentionCandidatesForChannel(people as never, "chn_a").map((person) => person.handle);
  assert.deepEqual(handles, ["alex", "kunde-a-am", "frei"]);
});

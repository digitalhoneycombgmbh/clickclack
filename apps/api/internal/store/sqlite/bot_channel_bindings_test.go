package sqlite

import (
	"context"
	"testing"

	"github.com/openclaw/clickclack/apps/api/internal/store"
)

// DHC: Ein gebundener Bot sieht und beschreibt nur seinen Channel; ein ungebundener Bot alle.
func TestBotChannelBindingRestrictsBotToItsChannel(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	owner, err := st.EnsureBootstrap(ctx, "Owner", "bind@example.com")
	if err != nil {
		t.Fatal(err)
	}
	workspace, err := st.CreateWorkspace(ctx, store.CreateWorkspaceInput{Name: "Fulfillment"}, owner.ID)
	if err != nil {
		t.Fatal(err)
	}
	own, _, err := st.CreateChannel(ctx, store.CreateChannelInput{WorkspaceID: workspace.ID, Name: "kunde-a", UserID: owner.ID})
	if err != nil {
		t.Fatal(err)
	}
	other, _, err := st.CreateChannel(ctx, store.CreateChannelInput{WorkspaceID: workspace.ID, Name: "kunde-b", UserID: owner.ID})
	if err != nil {
		t.Fatal(err)
	}
	newBot := func(handle string) store.User {
		bot, _, err := st.CreateBot(ctx, store.CreateBotInput{WorkspaceID: workspace.ID, DisplayName: handle, Handle: handle, TokenName: "t", Scopes: []string{"bot:write"}, CreatedBy: owner.ID})
		if err != nil {
			t.Fatal(err)
		}
		return bot
	}
	bound, free := newBot("kunde-a-am"), newBot("frei")
	if err := st.BindBotChannel(ctx, "kunde-a-am", own.ID); err != nil {
		t.Fatal(err)
	}
	if err := st.BindBotChannel(ctx, owner.ID, own.ID); err == nil {
		t.Fatal("humans must not be bindable")
	}

	if _, err := st.GetChannel(ctx, other.ID, bound.ID); err == nil {
		t.Fatal("bound bot can read a foreign channel")
	}
	if _, _, err := st.CreateMessage(ctx, store.CreateMessageInput{ChannelID: other.ID, AuthorID: bound.ID, Body: "hi"}); err == nil {
		t.Fatal("bound bot can post into a foreign channel")
	}
	if _, _, err := st.CreateMessage(ctx, store.CreateMessageInput{ChannelID: own.ID, AuthorID: bound.ID, Body: "hi"}); err != nil {
		t.Fatalf("bound bot cannot post into its channel: %v", err)
	}
	channels, err := st.ListChannels(ctx, workspace.ID, bound.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 || channels[0].ID != own.ID {
		t.Fatalf("bound bot lists %d channels, want only its own", len(channels))
	}
	if _, _, err := st.CreateMessage(ctx, store.CreateMessageInput{ChannelID: other.ID, AuthorID: free.ID, Body: "hi"}); err != nil {
		t.Fatalf("unbound bot lost access: %v", err)
	}
}

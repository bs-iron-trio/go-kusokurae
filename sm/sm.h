#ifndef BS_KUSOKURAE_SM_H
#define BS_KUSOKURAE_SM_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

#define KUSOKURAE_DECK_SIZE         33
#define KUSOKURAE_MAX_HAND_CARDS    22
#define KUSOKURAE_MAX_PLAYERS       4

struct kusokurae_game_state_t; // Forward declaration

typedef void (*state_transition_cb)(struct kusokurae_game_state_t *self, int32_t newstate, void *userdata);

typedef struct {
    int32_t np; // Number of players (3 or 4)
} kusokurae_game_config_t;

typedef struct {
    // State transition callback - to be called BEFORE each state change and end
    // of each round.
    void *userdata_of_state_transition;
    state_transition_cb state_transition;
} kusokurae_game_callbacks_t;

typedef enum {
    // 0 - Zero value
    KUSOKURAE_STATUS_NULL,

    // 1 - Struct initialized
    KUSOKURAE_STATUS_INIT,

    // 2 - Game in progress
    KUSOKURAE_STATUS_PLAY,

    // 3 - Game finished (you can retrieve results and/or start a new game)
    KUSOKURAE_STATUS_FINISH,

    // Keep this line at the bottom
    KUSOKURAE_STATUS_MAX,
} kusokurae_game_status_t;

typedef enum {
    // lit. "Shit"
    KUSOKURAE_SUIT_XIANG = -1,

    // lit. "Fried bread stick"
    KUSOKURAE_SUIT_YOUTIAO = 0,

    // lit. "Stuffed bun"
    KUSOKURAE_SUIT_BAOZI = 1,

    // In all the suits above, the number -1, 0 & 1 equal their values in a game,
    // but the following OTHER card type should be treated specially.
    KUSOKURAE_SUIT_OTHER = 2,
} kusokurae_card_suit_t;

typedef struct {
    // Sequence in the new, unshuffled deck. Higher value precedes lower
    // e.g. The newbiest card, Ghost, has a display_order of 33.
    // 0 indicates invalid data (unfilled slot).
    // Should be filled during global initialization and copied afterwards.
    uint32_t display_order;

    // Declared above (kusokurae_card_suit_t)
    int32_t suit;

    // 0~10 for BAOZI
    // 0~9 for YOUTIAO and XIANG
    // 10 for OTHER
    int32_t rank;

    // Bits 0~6: round index (counting from 1) in which the card is played
    // Bit 7: whether the card could be played in the current round
    // Bits 8~31: reserved
    uint32_t flags;
} kusokurae_card_t;

typedef enum {
    KUSOKURAE_ROUND_WAITING,
    KUSOKURAE_ROUND_ACTIVE,
    KUSOKURAE_ROUND_DONE,
} kusokurae_round_status_t;

typedef struct {
    // 1~4 (0 for invalid)
    int32_t index;

    // 1 - active (playing), 2 - already played
    int32_t active;

    // 22 card slots (reserved for playing with 2 decks)
    kusokurae_card_t cards[KUSOKURAE_MAX_HAND_CARDS];

    // The number of valid cards in hand.
    // When a card is played, it is removed from hand and all following cards
    //   should be moved ahead to keep the array consecutive.
    int32_t ncards;

    // If the player wins a round, he/she takes all cards played in that round.
    // cards_taken will always be multiples of player count.
    int32_t cards_taken;

    // The score accumulated from cards_taken.
    int32_t score;

    // When you say a player is busted, it means he/she is forced to play
    // forbidden moves because no other card's available.
    int32_t busted;
} kusokurae_player_t;

typedef enum {
    KUSOKURAE_SUCCESS,
    KUSOKURAE_ERROR_NULLPTR,
    KUSOKURAE_ERROR_BAD_NUMBER_OF_PLAYERS,
    KUSOKURAE_ERROR_UNINITIALIZED,
    KUSOKURAE_ERROR_NOT_IN_GAME,
    KUSOKURAE_ERROR_BUG_NOBODY_ACTIVE,
    KUSOKURAE_ERROR_CARD_NOT_FOUND,
    KUSOKURAE_ERROR_FORBIDDEN_MOVE,

    KUSOKURAE_ERROR_UNIMPLEMENTED,
    KUSOKURAE_ERROR_UNSPECIFIED,
} kusokurae_error_t;

typedef struct kusokurae_game_state_t {
    kusokurae_game_config_t cfg;
    int32_t status;

    // Max 4 players
    kusokurae_player_t players[KUSOKURAE_MAX_PLAYERS];

    // Finished round count
    int32_t nround;

    // Who has the ghost in hand?
    int32_t ghost_holder_index;

    // Rank leader in the current round.
    // Set to -1 before anyone plays and updated on each play.
    int32_t high_ranker_index;

    // Cards played in the current round.
    // players[n]'s move is placed in current_round[n].
    kusokurae_card_t current_round[KUSOKURAE_MAX_PLAYERS];

    // 8 bytes of state for random number generator.
    uint64_t rng_state;

    // Game-specific callbacks should be put at the bottom, because their sizes
    // are machine-dependent.
    kusokurae_game_callbacks_t cbs;
} kusokurae_game_state_t;

typedef struct {
    // On screen: "Round <seq>"
    int32_t seq;

    // Whether there is a ghost
    int32_t is_doubled;

    // Total score in cards played
    int32_t score_on_board;

    // The current winning player
    int32_t round_winner;

    // Moves made in this round, ordered chronologically (e.g. if there're 3
    // players and the trick leader is 2P, then moves[0] is 2P's move, moves[1]
    // is 3P's move, moves[2] is 1P's move, and moves[3] is unused)
    kusokurae_card_t moves[KUSOKURAE_MAX_PLAYERS];
} kusokurae_round_state_t;

void kusokurae_global_init();

void kusokurae_set_prng(int16_t (*fn)(void *));

kusokurae_error_t kusokurae_game_init(kusokurae_game_state_t *self,
                                      kusokurae_game_config_t *cfg,
                                      kusokurae_game_callbacks_t *cbs);

kusokurae_error_t kusokurae_game_start(kusokurae_game_state_t *self);

kusokurae_error_t kusokurae_game_play(kusokurae_game_state_t *self,
                                      kusokurae_card_t card);

int kusokurae_game_is_final_round(kusokurae_game_state_t *self);

kusokurae_player_t *kusokurae_get_active_player(kusokurae_game_state_t *self);

void kusokurae_get_round_state(kusokurae_game_state_t *self,
                               kusokurae_round_state_t *out);

int kusokurae_card_is_playable(kusokurae_card_t card);

int kusokurae_card_round_played(kusokurae_card_t card);

#ifdef __cplusplus
}
#endif

#endif // BS_KUSOKURAE_SM_H

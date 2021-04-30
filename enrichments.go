package mirvpgl

// ( see https://wiki.alliedmods.net/Counter-Strike:_Global_Offensive_Events )
var (
	enrichments = Enrichments{
		"player_death": {
			"userid":   &UserIDEnrichment{},
			"attacker": &UserIDEnrichment{},
			"assister": &UserIDEnrichment{},
		},
		"other_death": {
			"attacker": &UserIDEnrichment{},
		},
		"player_hurt": {
			"userid":   &UserIDEnrichment{},
			"attacker": &UserIDEnrichment{},
		},
		"item_purchase": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_beginplant": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_abortplant": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_planted": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_defused": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_exploded": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_pickup": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_dropped": {
			"userid":   &UserIDEnrichment{},
			"entindex": &EntityNumEnrichment{},
		},
		"defuser_dropped": {
			"entityid": &EntityNumEnrichment{},
		},
		"defuser_pickup": {
			"entityid": &EntityNumEnrichment{},
			"userid":   &UserIDEnrichment{},
		},
		"bomb_begindefuse": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_abortdefuse": {
			"userid": &UserIDEnrichment{},
		},
		"hostage_follows": {
			"userid":  &UserIDEnrichment{},
			"hostage": &EntityNumEnrichment{},
		},
		"hostage_hurt": {
			"userid":  &UserIDEnrichment{},
			"hostage": &EntityNumEnrichment{},
		},
		"hostage_killed": {
			"userid":  &UserIDEnrichment{},
			"hostage": &EntityNumEnrichment{},
		},
		"hostage_rescued": {
			"userid":  &UserIDEnrichment{},
			"hostage": &EntityNumEnrichment{},
		},
		"hostage_stops_following": {
			"userid":  &UserIDEnrichment{},
			"hostage": &EntityNumEnrichment{},
		},
		"hostage_call_for_help": {
			"hostage": &EntityNumEnrichment{},
		},
		"vip_escaped": {
			"userid": &UserIDEnrichment{},
		},
		"player_radio": {
			"userid": &UserIDEnrichment{},
		},
		"bomb_beep": {
			"entindex": &EntityNumEnrichment{},
		},
		"weapon_fire": {
			"userid": &UserIDEnrichment{},
		},
		"weapon_fire_on_empty": {
			"userid": &UserIDEnrichment{},
		},
		"grenade_thrown": {
			"userid": &UserIDEnrichment{},
		},
		"weapon_outofammo": {
			"userid": &UserIDEnrichment{},
		},
		"weapon_reload": {
			"userid": &UserIDEnrichment{},
		},
		"weapon_zoom": {
			"userid": &UserIDEnrichment{},
		},
		"silencer_detach": {
			"userid": &UserIDEnrichment{},
		},
		"inspect_weapon": {
			"userid": &UserIDEnrichment{},
		},
		"weapon_zoom_rifle": {
			"userid": &UserIDEnrichment{},
		},
		"player_spawned": {
			"userid": &UserIDEnrichment{},
		},
		"item_pickup": {
			"userid": &UserIDEnrichment{},
		},
		"item_pickup_failed": {
			"userid": &UserIDEnrichment{},
		},
		"item_remove": {
			"userid": &UserIDEnrichment{},
		},
		"ammo_pickup": {
			"userid": &UserIDEnrichment{},
			"index":  &EntityNumEnrichment{},
		},
		"item_equip": {
			"userid": &UserIDEnrichment{},
		},
		"enter_buyzone": {
			"userid": &UserIDEnrichment{},
		},
		"exit_buyzone": {
			"userid": &UserIDEnrichment{},
		},
		"enter_bombzone": {
			"userid": &UserIDEnrichment{},
		},
		"exit_bombzone": {
			"userid": &UserIDEnrichment{},
		},
		"enter_rescue_zone": {
			"userid": &UserIDEnrichment{},
		},
		"exit_rescue_zone": {
			"userid": &UserIDEnrichment{},
		},
		"silencer_off": {
			"userid": &UserIDEnrichment{},
		},
		"silencer_on": {
			"userid": &UserIDEnrichment{},
		},
		"buymenu_open": {
			"userid": &UserIDEnrichment{},
		},
		"buymenu_close": {
			"userid": &UserIDEnrichment{},
		},
		"round_end": {
			"winner": &UserIDEnrichment{},
		},
		"grenade_bounce": {
			"userid": &UserIDEnrichment{},
		},
		"hegrenade_detonate": {
			"userid": &UserIDEnrichment{},
		},
		"flashbang_detonate": {
			"userid": &UserIDEnrichment{},
		},
		"smokegrenade_detonate": {
			"userid": &UserIDEnrichment{},
		},
		"smokegrenade_expired": {
			"userid": &UserIDEnrichment{},
		},
		"molotov_detonate": {
			"userid": &UserIDEnrichment{},
		},
		"decoy_detonate": {
			"userid": &UserIDEnrichment{},
		},
		"decoy_started": {
			"userid": &UserIDEnrichment{},
		},
		"tagrenade_detonate": {
			"userid": &UserIDEnrichment{},
		},
		"decoy_firing": {
			"userid": &UserIDEnrichment{},
		},
		"bullet_impact": {
			"userid": &UserIDEnrichment{},
		},
		"player_footstep": {
			"userid": &UserIDEnrichment{},
		},
		"player_jump": {
			"userid": &UserIDEnrichment{},
		},
		"player_blind": {
			"userid":   &UserIDEnrichment{},
			"entityid": &EntityNumEnrichment{},
		},
		"player_falldamage": {
			"userid": &UserIDEnrichment{},
		},
		"door_moving": {
			"entityid": &EntityNumEnrichment{},
			"userid":   &UserIDEnrichment{},
		},
		"spec_target_updated": {
			"userid": &UserIDEnrichment{},
		},
		"player_avenged_teammate": {
			"avenger_id":        &UserIDEnrichment{},
			"avenged_player_id": &UserIDEnrichment{},
		},
		"round_mvp": {
			"userid": &UserIDEnrichment{},
		},
		"player_decal": {
			"userid": &UserIDEnrichment{},
		},

		// ... left out the gg / gungame shit, feel free to add it ...

		"player_reset_vote": {
			"userid": &UserIDEnrichment{},
		},
		"start_vote": {
			"userid": &UserIDEnrichment{},
		},
		"player_given_c4": {
			"userid": &UserIDEnrichment{},
		},
		"player_become_ghost": {
			"userid": &UserIDEnrichment{},
		},

		// ... left out the tr shit, feel free to add it ...

		"jointeam_failed": {
			"userid": &UserIDEnrichment{},
		},
		"teamchange_pending": {
			"userid": &UserIDEnrichment{},
		},
		"ammo_refill": {
			"userid": &UserIDEnrichment{},
		},

		// ... left out the dangerzone shit, feel free to add it ...

		// others:

		"weaponhud_selection": {
			"userid": &UserIDEnrichment{},
		},
	}
)

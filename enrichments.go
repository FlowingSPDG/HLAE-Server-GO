package mirvpgl

// ( see https://wiki.alliedmods.net/Counter-Strike:_Global_Offensive_Events )
var (
	enrichments = Enrichments{
		"player_death": {
			"userid":   newUserIDEnrichment(),
			"attacker": newUserIDEnrichment(),
			"assister": newUserIDEnrichment(),
		},
		"other_death": {
			"attacker": newUserIDEnrichment(),
		},
		"player_hurt": {
			"userid":   newUserIDEnrichment(),
			"attacker": newUserIDEnrichment(),
		},
		"item_purchase": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_beginplant": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_abortplant": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_planted": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_defused": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_exploded": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_pickup": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_dropped": {
			"userid":   newUserIDEnrichment(),
			"entindex": newEntityNumEnrichment(),
		},
		"defuser_dropped": {
			"entityid": newEntityNumEnrichment(),
		},
		"defuser_pickup": {
			"entityid": newEntityNumEnrichment(),
			"userid":   newUserIDEnrichment(),
		},
		"bomb_begindefuse": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_abortdefuse": {
			"userid": newUserIDEnrichment(),
		},
		"hostage_follows": {
			"userid":  newUserIDEnrichment(),
			"hostage": newEntityNumEnrichment(),
		},
		"hostage_hurt": {
			"userid":  newUserIDEnrichment(),
			"hostage": newEntityNumEnrichment(),
		},
		"hostage_killed": {
			"userid":  newUserIDEnrichment(),
			"hostage": newEntityNumEnrichment(),
		},
		"hostage_rescued": {
			"userid":  newUserIDEnrichment(),
			"hostage": newEntityNumEnrichment(),
		},
		"hostage_stops_following": {
			"userid":  newUserIDEnrichment(),
			"hostage": newEntityNumEnrichment(),
		},
		"hostage_call_for_help": {
			"hostage": newEntityNumEnrichment(),
		},
		"vip_escaped": {
			"userid": newUserIDEnrichment(),
		},
		"player_radio": {
			"userid": newUserIDEnrichment(),
		},
		"bomb_beep": {
			"entindex": newEntityNumEnrichment(),
		},
		"weapon_fire": {
			"userid": newUserIDEnrichment(),
		},
		"weapon_fire_on_empty": {
			"userid": newUserIDEnrichment(),
		},
		"grenade_thrown": {
			"userid": newUserIDEnrichment(),
		},
		"weapon_outofammo": {
			"userid": newUserIDEnrichment(),
		},
		"weapon_reload": {
			"userid": newUserIDEnrichment(),
		},
		"weapon_zoom": {
			"userid": newUserIDEnrichment(),
		},
		"silencer_detach": {
			"userid": newUserIDEnrichment(),
		},
		"inspect_weapon": {
			"userid": newUserIDEnrichment(),
		},
		"weapon_zoom_rifle": {
			"userid": newUserIDEnrichment(),
		},
		"player_spawned": {
			"userid": newUserIDEnrichment(),
		},
		"item_pickup": {
			"userid": newUserIDEnrichment(),
		},
		"item_pickup_failed": {
			"userid": newUserIDEnrichment(),
		},
		"item_remove": {
			"userid": newUserIDEnrichment(),
		},
		"ammo_pickup": {
			"userid": newUserIDEnrichment(),
			"index":  newEntityNumEnrichment(),
		},
		"item_equip": {
			"userid": newUserIDEnrichment(),
		},
		"enter_buyzone": {
			"userid": newUserIDEnrichment(),
		},
		"exit_buyzone": {
			"userid": newUserIDEnrichment(),
		},
		"enter_bombzone": {
			"userid": newUserIDEnrichment(),
		},
		"exit_bombzone": {
			"userid": newUserIDEnrichment(),
		},
		"enter_rescue_zone": {
			"userid": newUserIDEnrichment(),
		},
		"exit_rescue_zone": {
			"userid": newUserIDEnrichment(),
		},
		"silencer_off": {
			"userid": newUserIDEnrichment(),
		},
		"silencer_on": {
			"userid": newUserIDEnrichment(),
		},
		"buymenu_open": {
			"userid": newUserIDEnrichment(),
		},
		"buymenu_close": {
			"userid": newUserIDEnrichment(),
		},
		"round_end": {
			"winner": newUserIDEnrichment(),
		},
		"grenade_bounce": {
			"userid": newUserIDEnrichment(),
		},
		"hegrenade_detonate": {
			"userid": newUserIDEnrichment(),
		},
		"flashbang_detonate": {
			"userid": newUserIDEnrichment(),
		},
		"smokegrenade_detonate": {
			"userid": newUserIDEnrichment(),
		},
		"smokegrenade_expired": {
			"userid": newUserIDEnrichment(),
		},
		"molotov_detonate": {
			"userid": newUserIDEnrichment(),
		},
		"decoy_detonate": {
			"userid": newUserIDEnrichment(),
		},
		"decoy_started": {
			"userid": newUserIDEnrichment(),
		},
		"tagrenade_detonate": {
			"userid": newUserIDEnrichment(),
		},
		"decoy_firing": {
			"userid": newUserIDEnrichment(),
		},
		"bullet_impact": {
			"userid": newUserIDEnrichment(),
		},
		"player_footstep": {
			"userid": newUserIDEnrichment(),
		},
		"player_jump": {
			"userid": newUserIDEnrichment(),
		},
		"player_blind": {
			"userid":   newUserIDEnrichment(),
			"entityid": newEntityNumEnrichment(),
		},
		"player_falldamage": {
			"userid": newUserIDEnrichment(),
		},
		"door_moving": {
			"entityid": newEntityNumEnrichment(),
			"userid":   newUserIDEnrichment(),
		},
		"spec_target_updated": {
			"userid": newUserIDEnrichment(),
		},
		"player_avenged_teammate": {
			"avenger_id":        newUserIDEnrichment(),
			"avenged_player_id": newUserIDEnrichment(),
		},
		"round_mvp": {
			"userid": newUserIDEnrichment(),
		},
		"player_decal": {
			"userid": newUserIDEnrichment(),
		},

		// ... left out the gg / gungame shit, feel free to add it ...

		"player_reset_vote": {
			"userid": newUserIDEnrichment(),
		},
		"start_vote": {
			"userid": newUserIDEnrichment(),
		},
		"player_given_c4": {
			"userid": newUserIDEnrichment(),
		},
		"player_become_ghost": {
			"userid": newUserIDEnrichment(),
		},

		// ... left out the tr shit, feel free to add it ...

		"jointeam_failed": {
			"userid": newUserIDEnrichment(),
		},
		"teamchange_pending": {
			"userid": newUserIDEnrichment(),
		},
		"ammo_refill": {
			"userid": newUserIDEnrichment(),
		},

		// ... left out the dangerzone shit, feel free to add it ...

		// others:

		"weaponhud_selection": {
			"userid": newUserIDEnrichment(),
		},
	}
)

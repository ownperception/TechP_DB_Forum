package middlefunc

import "strconv"

func ReqThreadsFlags(params map[string]string, flags *map[string]string) {
	if params["desc"] == "true" {
		(*flags)["sortflag"] = "DESC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = "where created <= '" + params["since"] + "' "
		}
	} else {
		(*flags)["sortflag"] = "ASC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = "where created >= '" + params["since"] + "' "
		}
	}

	if params["limit"] != "" {
		(*flags)["limitflag"] = "limit " + params["limit"]
	}
}

func ReqUsersFlags(params map[string]string, flags *map[string]string) {
	if params["desc"] == "true" {
		(*flags)["sortflag"] = "DESC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = "and author < '" + params["since"] + "'"
		}
	} else {
		(*flags)["sortflag"] = "ASC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = "and author > '" + params["since"] + "'"
		}
	}

	if params["limit"] != "" {
		(*flags)["limitflag"] = "limit " + params["limit"]
	}
}

func ReqFlatFlags(params map[string]string, flags *map[string]string) {
	if params["desc"] == "true" {
		(*flags)["sortflag"] = "DESC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = " and id  < " + params["since"]
		}
	} else {
		(*flags)["sortflag"] = "ASC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = " and id > " + params["since"]
		}
	}

	if params["limit"] != "" {
		(*flags)["limitflag"] = "limit " + params["limit"]
	}
}

func ReqTreeFlags(params map[string]string, flags *map[string]string) {

	if params["desc"] == "true" {
		(*flags)["sortflag"] = " order by path[0],path DESC "
		if params["since"] != "" {
			(*flags)["sinceflag"] = " where path < (select path from post_tree where id = " + params["since"] + " ) "
		}
	} else {
		(*flags)["sortflag"] = " order by path[0],path ASC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = " where path > (select path from post_tree where id = " + params["since"] + " ) "
		}
	}

	if params["limit"] != "" {
		(*flags)["limitflag"] = "limit " + params["limit"]
	}
}

func ReqParTreeFlags(params map[string]string, flags *map[string]string) {
	if params["desc"] == "true" {
		(*flags)["descflag"] = " desc "
		(*flags)["sortflag"] = "order by path[1] DESC,path"
		if params["since"] != "" {
			(*flags)["sinceflag"] = " where path[1] < (select path[1] from post_tree where id = " + params["since"] + " ) "
		}
	} else {
		(*flags)["descflag"] = " asc "
		(*flags)["sortflag"] = " order by path[1] ,path ASC"
		if params["since"] != "" {
			(*flags)["sinceflag"] = " where path[1] > (select path[1] from post_tree where id = " + params["since"] + " ) "
		}
	}

	if params["limit"] != "" {
		(*flags)["limitflag"] = " where r <= " + params["limit"]
	}
}

func FlagSlugOrId(id string) string {
	var id_flag string
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		id_flag = "slug = '" + id + "'"
	} else {
		id_flag = "id = '" + id + "'"
	}
	return id_flag
}

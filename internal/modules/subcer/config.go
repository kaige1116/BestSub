package subcer

const pref = `
common:
  api_mode: true
  base_path: base
  clash_rule_base: base/all_base.tpl
  surge_rule_base: base/all_base.tpl
  surfboard_rule_base: base/all_base.tpl
  mellow_rule_base: base/all_base.tpl
  quan_rule_base: base/all_base.tpl
  quanx_rule_base: base/all_base.tpl
  loon_rule_base: base/all_base.tpl
  sssub_rule_base: base/all_base.tpl
  singbox_rule_base: base/all_base.tpl
  reload_conf_on_request: false


node_pref:
  clash_use_new_field_name: true
  clash_proxies_style: flow
  clash_proxy_groups_style: block
  singbox_add_clash_modes: true


emojis:
  add_emoji: false
  remove_old_emoji: false

template:
  globals:
  - {key: clash.http_port, value: 7890}
  - {key: clash.socks_port, value: 7891}
  - {key: clash.allow_lan, value: true}
  - {key: clash.log_level, value: info}
  - {key: clash.external_controller, value: '127.0.0.1:9090'}
  - {key: singbox.allow_lan, value: true}
  - {key: singbox.mixed_port, value: 2080}

server:
  listen: %s
  port: %d

advanced:
  log_level: debug
  max_pending_connections: 10240
  max_concurrent_threads: 100
  enable_cache: true
  cache_subscription: 0
  cache_config: 86400
  cache_ruleset: 86400
`

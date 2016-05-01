package apparmor

// baseTemplate for AppArmor profiles
const baseTemplate = `
{{range $value := .Imports}}{{$value}}
{{end}}

profile {{.Name}} flags=(attach_disconnected,mediate_deleted) {
{{range $value := .InnerImports}}  {{$value}}
{{end}}

{{if .Network.Protocols}}
{{range $value := .Network.Protocols}}  network inet {{$value}},
{{end}}{{else}}
  network,
{{end}}
{{if .Network.Raw}}{{else}}  deny network raw,
{{end}}
{{if .Network.Packet}}{{else}}  deny network packet,
{{end}}
  file,
  umount,

{{range $value := .Filesystem.ReadOnlyPaths}}  deny {{$value}} wl,
{{end}}
{{range $value := .Filesystem.LogOnWritePaths}}  audit {{$value}} w,
{{end}}
{{range $value := .Filesystem.WritablePaths}}  {{$value}} w,
{{end}}
{{range $value := .Filesystem.AllowExec}}  {{$value}} ix,
{{end}}
{{range $value := .Filesystem.DenyExec}}  deny {{$value}} mrwklx,
{{end}}
{{if .Capabilities.Deny}}
  capability,
{{range $value := .Capabilities.Deny}}  deny capability {{$value}},
{{end}}{{else}}
{{range $value := .Capabilities.Allow}}  capability {{$value}},
{{end}}{{end}}
  deny @{PROC}/* w,   # deny write for all files directly in /proc (not in a subdir)
  deny @{PROC}/{[^1-9],[^1-9][^0-9],[^1-9s][^0-9y][^0-9s],[^1-9][^0-9][^0-9][^0-9]*}/** w,
  deny @{PROC}/sys/[^k]** w,  # deny /proc/sys except /proc/sys/k* (effectively /proc/sys/kernel)
  deny @{PROC}/sys/kernel/{?,??,[^s][^h][^m]**} w,  # deny everything except shm* in /proc/sys/kernel/
  deny @{PROC}/sysrq-trigger rwklx,
  deny @{PROC}/mem rwklx,
  deny @{PROC}/kmem rwklx,
  deny @{PROC}/kcore rwklx,
  deny mount,
  deny /sys/[^f]*/** wklx,
  deny /sys/f[^s]*/** wklx,
  deny /sys/fs/[^c]*/** wklx,
  deny /sys/fs/c[^g]*/** wklx,
  deny /sys/fs/cg[^r]*/** wklx,
  deny /sys/firmware/efi/efivars/** rwklx,
  deny /sys/kernel/security/** rwklx,
}
`

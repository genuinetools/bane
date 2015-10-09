package apparmor

// baseTemplate for AppArmor profiles
const baseTemplate = `
{{range $value := .Imports}}{{$value}}
{{end}}

profile {{.Name}} flags=(attach_disconnected,mediate_deleted) {
{{range $value := .InnerImports}}  {{$value}}
{{end}}

  network,
{{if .Network.Raw}}{{else}}  deny network raw,
{{end}}
{{if .Network.Packet}}{{else}}  deny network packet,
{{end}}
  capability,
  file,
  umount,

{{range $value := .ReadOnlyPaths}}  deny {{$value}} wl,
{{end}}
{{range $value := .LogOnWritePaths}}  audit {{$value}} w,
{{end}}
{{range $value := .WritablePaths}}  {{$value}} w,
{{end}}
{{range $value := .Executables.Allow}}  {{$value}} ix,
{{end}}
{{range $value := .Executables.Deny}}  deny {{$value}} mrwklx,
{{end}}

  deny @{PROC}/{*,**^[0-9*],sys/kernel/shm*} wkx,
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

